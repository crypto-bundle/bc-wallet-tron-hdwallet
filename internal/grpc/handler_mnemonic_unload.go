package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
)

const (
	MethodNameUnLoadMnemonic = "UnLoadMnemonic"
)

type unLoadMnemonicHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *unLoadMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &UnLoadMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	err = h.walletPoolSvc.UnloadWalletUnit(tCtx, vf.WalletUUIDRaw)
	if err != nil {
		h.l.Error("unable to unload mnemonic wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pbApi.UnLoadMnemonicResponse{
		MnemonicIdentity: req.MnemonicIdentity,
	}, nil
}

func MakeUnLoadMnemonicHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *unLoadMnemonicHandler {
	return &unLoadMnemonicHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodNameUnLoadMnemonic)),
		walletPoolSvc: walletPoolSvc,
	}
}
