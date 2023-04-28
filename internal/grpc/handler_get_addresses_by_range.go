package grpc

import (
	"context"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/types"
	"sync"

	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/app"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/forms"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetDerivationAddressByRange = "GetDerivationAddressByRange"
)

type GetDerivationAddressByRangeHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
	respPool      sync.Pool
}

type respAddrList []*pbApi.DerivationAddressIdentity

func (l *respAddrList) Reset() {

}

// nolint:funlen // fixme
func (h *GetDerivationAddressByRangeHandler) Handle(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &forms.DerivationAddressByRangeForm{}
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

	return h.processRequest(ctx, vf, walletPubData, mnemoWalletData)
}

func (h *GetDerivationAddressByRangeHandler) processRequest(ctx context.Context,
	vf *forms.DerivationAddressByRangeForm,
	walletPubData *types.PublicWalletData,
	mnemoWalletData *types.PublicMnemonicWalletData,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	var err error

	rangeSize := (vf.AddressIndexTo - vf.AddressIndexFrom) + 1
	filedData := make([]*pbApi.DerivationAddressIdentity, rangeSize)

	marshallerCallback := func(addressIdx, position uint32, address string) {
		addressEntity := h.respPool.Get().(*pbApi.DerivationAddressIdentity)

		addressEntity.AccountIndex = vf.AccountIndex
		addressEntity.InternalIndex = vf.AccountIndex
		addressEntity.AddressIndex = addressIdx
		addressEntity.Address = address

		filedData[position] = addressEntity
		return
	}

	err = h.walletSrv.GetAddressesByPathByRange(ctx, vf.WalletUUIDRaw, vf.MnemonicWalletUUIDRaw,
		vf.AccountIndex, vf.InternalIndex,
		vf.AddressIndexFrom, vf.AddressIndexTo, marshallerCallback)
	if err != nil {
		h.l.Error("unable get derivative addresses by range", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	response, err := h.marshallerSrv.MarshallGetAddressByRange(walletPubData, mnemoWalletData, filedData)
	if err != nil {
		h.l.Error("unable to marshall get addresses data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer func(clearedSize uint32) {
		go func(size uint32) {
			for i := uint32(0); i != size; i++ {
				h.respPool.Put(filedData[i])
			}
		}(clearedSize)
	}(rangeSize)

	return response, nil
}

func MakeGetDerivationAddressByRangeHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *GetDerivationAddressByRangeHandler {
	return &GetDerivationAddressByRangeHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodGetDerivationAddressByRange)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
		respPool: sync.Pool{New: func() any {
			return new(pbApi.DerivationAddressIdentity)
		}},
	}
}
