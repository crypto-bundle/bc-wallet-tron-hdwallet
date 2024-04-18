package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"
	"time"

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

type mnemonicGeneratorService interface {
	Generate(ctx context.Context) (string, error)
}

type mnemonicValidatorService interface {
	IsMnemonicValid(mnemonic string) bool
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewTronWallet(account, change, address uint32) (*hdwallet.Tron, error)

	ClearSecrets()
}

type walletPoolService interface {
	AddAndStartWalletUnit(_ context.Context,
		mnemonicWalletUUID uuid.UUID,
		timeToLive time.Duration,
		mnemonicEncryptedData []byte,
	) error
	LoadAddressByPath(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*string, error)
	UnloadWalletUnit(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
	) (*uuid.UUID, error)
	GetAddressByPath(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*string, error)
	GetAddressesByPathByRange(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*string, []byte, error)
}

type generateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.GenerateMnemonicRequest) (*pbApi.GenerateMnemonicResponse, error)
}

type validateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.ValidateMnemonicRequest) (*pbApi.ValidateMnemonicResponse, error)
}

type loadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.LoadMnemonicRequest) (*pbApi.LoadMnemonicResponse, error)
}

type unLoadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.UnLoadMnemonicRequest) (*pbApi.UnLoadMnemonicResponse, error)
}

type unLoadMultipleMnemonicsHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.UnLoadMultipleMnemonicsRequest,
	) (*pbApi.UnLoadMultipleMnemonicsResponse, error)
}

type encryptMnemonicHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.EncryptMnemonicRequest,
	) (*pbApi.EncryptMnemonicResponse, error)
}

type getDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressByRangeRequest,
	) (*pbApi.DerivationAddressByRangeResponse, error)
}
type loadDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.LoadDerivationAddressRequest,
	) (*pbApi.LoadDerivationAddressResponse, error)
}

type signDataHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.SignDataRequest,
	) (*pbApi.SignDataResponse, error)
}

type getDerivationAddressHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressRequest,
	) (*pbApi.DerivationAddressResponse, error)
}
