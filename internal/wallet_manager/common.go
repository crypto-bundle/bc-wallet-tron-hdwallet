package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
)

type configService interface {
	GetDefaultWalletUnloadInterval() time.Duration
}

type WalletPoolUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	LoadAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (string, error)
	GetAddressByPath(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		account, change, index uint32,
		transactionData []byte,
	) (*string, []byte, error)
}

type walletPoolService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	SetWalletUnits(ctx context.Context,
		walletUnits map[uuid.UUID]WalletPoolUnitService,
	) error
	AddAndStartWalletUnit(ctx context.Context,
		walletUUID uuid.UUID,
		walletUnit WalletPoolUnitService,
	) error
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignTransaction(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*types.PublicSignTxData, error)
}

type mnemonicWalletConfig interface {
	GetMnemonicWalletPurpose() string
	GetMnemonicWalletHash() string
	IsHotWallet() bool
}

type walleter interface {
	GetAddress() (string, error)
	GetPubKey() string
	GetPrvKey() (string, error)
	GetPath() string
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewTronWallet(account, change, address uint32) (*hdwallet.Tron, error)

	ClearSecrets()
}

type mnemonicGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}
