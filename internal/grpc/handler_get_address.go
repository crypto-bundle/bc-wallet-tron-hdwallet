package grpc

import (
	"context"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/app"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/forms"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddress = "GetDerivationAddress"
)

type GetDerivationAddressHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *GetDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &forms.GetDerivationAddressForm{}
	valid, err := vf.LoadAndValidate(ctx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	walletPubData, err := h.walletSrv.GetWalletByUUID(ctx, vf.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletPubData == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	mnemoWalletData, isExists := walletPubData.MnemonicWalletsByUUID[vf.MnemonicWalletUUIDRaw]
	if !isExists {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	addressData, err := h.walletSrv.GetAddressByPath(ctx, vf.WalletUUIDRaw, vf.MnemonicWalletUUIDRaw,
		vf.AccountIndex, vf.InternalIndex, vf.AddressIndex)
	if err != nil {
		return nil, err
	}

	marshalledData, err := h.marshallerSrv.MarshallGetAddressData(walletPubData, mnemoWalletData, addressData)
	if err != nil {
		h.l.Error("unable to marshall public address data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return marshalledData, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetDerivationAddressHandler {
	return &GetDerivationAddressHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddress)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
