package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
)

const (
	MethodNameLoadMnemonic = "LoadMnemonic"
)

type loadMnemonicHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *loadMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.LoadMnemonicRequest,
) (*pbApi.LoadMnemonicResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &LoadMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	err = h.walletPoolSvc.AddAndStartWalletUnit(tCtx, vf.WalletUUIDRaw,
		time.Duration(vf.TimeToLive),
		vf.EncryptedMnemonicData)
	if err != nil {
		h.l.Error("unable load mnemonic wallet", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pbApi.LoadMnemonicResponse{
		MnemonicIdentity: req.MnemonicIdentity,
	}, nil
}

func MakeLoadMnemonicHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *loadMnemonicHandler {
	return &loadMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameLoadMnemonic)),

		walletPoolSvc: walletPoolSvc,
	}
}
