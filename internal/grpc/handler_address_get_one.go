package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameGetDerivationAddress = "GetDerivationAddress"
)

type getDerivationAddressHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *getDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (*pbApi.DerivationAddressResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &AddressForm{}
	valid, err := vf.LoadAndValidateGetAddrReq(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	addr, err := h.walletPoolSvc.GetAddressByPath(tCtx, vf.WalletUUIDRaw,
		vf.AccountIndex,
		vf.InternalIndex,
		vf.AddressIndex)
	if err != nil {
		h.l.Error("unable to get address by path", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if addr == nil {
		return nil, status.Error(codes.ResourceExhausted, "wallet not loaded")
	}

	req.AddressIdentity.Address = *addr

	return &pbApi.DerivationAddressResponse{
		MnemonicWalletIdentifier: req.MnemonicWalletIdentifier,
		AddressIdentity:          req.AddressIdentity,
	}, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *getDerivationAddressHandler {
	return &getDerivationAddressHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGetDerivationAddress)),

		walletPoolSvc: walletPoolSvc,
	}
}
