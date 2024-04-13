package grpc_hdwallet

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"go.uber.org/zap"
)

const (
	MethodNameUnLoadMultipleMnemonics = "UnLoadMultipleMnemonics"
)

type unLoadMultipleMnemonicsHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *unLoadMultipleMnemonicsHandler) Handle(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (*pbApi.UnLoadMultipleMnemonicsResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	return nil, nil
}

func MakeUnLoadMultipleMnemonicsHandler(loggerEntry *zap.Logger) *unLoadMultipleMnemonicsHandler {
	return &unLoadMultipleMnemonicsHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameUnLoadMultipleMnemonics)),
	}
}
