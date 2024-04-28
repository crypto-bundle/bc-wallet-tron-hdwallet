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
	MethodNameGetDerivationsAddresses = "GetDerivationsAddresses"
)

type getDerivationsAddressesHandler struct {
	l *zap.Logger

	walletPoolSvc walletPoolService
}

// nolint:funlen // fixme
func (h *getDerivationsAddressesHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &derivationAddressByRangeForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	// marshallerSvc - instance of marshaller with closure usage
	marshallerSvc := newAddrRangeMarshaller(vf.RangeSize)

	err = h.walletPoolSvc.GetAddressesByPathByRange(tCtx, vf.MnemonicWalletUUIDRaw,
		vf, marshallerSvc.MarshallPath)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	return &pbApi.DerivationAddressByRangeResponse{
		MnemonicWalletIdentifier: req.MnemonicWalletIdentifier,
		AddressIdentitiesCount:   uint64(vf.RangeSize),
		AddressIdentities:        marshallerSvc.GetMarshaled(),
	}, nil
}

func MakeGetDerivationsAddressesHandler(loggerEntry *zap.Logger,
	walletPoolSvc walletPoolService,
) *getDerivationsAddressesHandler {
	return &getDerivationsAddressesHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGetDerivationsAddresses)),

		walletPoolSvc: walletPoolSvc,
	}
}
