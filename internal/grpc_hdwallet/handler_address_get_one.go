package grpc_hdwallet

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
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
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	return nil, nil
}

func MakeGetDerivationAddressHandler(loggerEntry *zap.Logger,
) *getDerivationAddressHandler {
	return &getDerivationAddressHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGetDerivationAddress)),
	}
}
