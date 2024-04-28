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
	MethodNameLoadDerivationAddress = "LoadDerivationAddress"
)

type loadDerivationAddressHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *loadDerivationAddressHandler) Handle(ctx context.Context,
	req *pbApi.LoadDerivationAddressRequest,
) (*pbApi.LoadDerivationAddressResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &AddressForm{}
	valid, err := vf.LoadAndValidateLoadAddrReq(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	addr, err := h.walletPoolSvc.LoadAddressByPath(tCtx, vf.WalletUUIDRaw,
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

	req.AddressIdentifier.Address = *addr

	return &pbApi.LoadDerivationAddressResponse{
		MnemonicWalletIdentifier: req.MnemonicWalletIdentifier,
		TxOwnerIdentity:          req.AddressIdentifier,
	}, nil
}

func MakeLoadDerivationAddressHandlerHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *loadDerivationAddressHandler {
	return &loadDerivationAddressHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameLoadDerivationAddress)),
		walletPoolSvc: walletPoolSvc,
	}
}
