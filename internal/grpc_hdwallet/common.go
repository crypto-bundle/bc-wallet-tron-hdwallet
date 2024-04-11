package grpc_hdwallet

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
)

const (
	MethodNameTag = "method_name"
)

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetConnectionPath() string
}

type walletManagerService interface {
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*string, error)
	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*string, []byte, error)
}

type generateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.GenerateMnemonicRequest) (*pbApi.GenerateMnemonicResponse, error)
}

type loadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.LoadMnemonicRequest) (*pbApi.LoadMnemonicResponse, error)
}

type unLoadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.UnLoadMnemonicRequest) (*pbApi.UnLoadMnemonicResponse, error)
}

type getDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressByRangeRequest,
	) (*pbApi.DerivationAddressByRangeResponse, error)
}

type signDataHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.SignTransactionRequest,
	) (*pbApi.SignTransactionResponse, error)
}

type getDerivationAddressHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressRequest,
	) (*pbApi.DerivationAddressResponse, error)
}
