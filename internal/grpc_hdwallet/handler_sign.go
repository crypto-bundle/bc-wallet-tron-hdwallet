package grpc_hdwallet

import (
	"context"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"go.uber.org/zap"
)

const (
	MethodNameSignMnemonic = "Sign"
)

type signDataHandler struct {
	l *zap.Logger
}

// nolint:funlen // fixme
func (h *signDataHandler) Handle(ctx context.Context,
	req *pbApi.SignDataRequest,
) (*pbApi.SignDataResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	return nil, nil
}

func MakeSignDataHandler(loggerEntry *zap.Logger,
) *signDataHandler {
	return &signDataHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameSignMnemonic)),
	}
}
