package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
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
	) (*string, error)
	GetAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		account, change, index uint32,
		dataForSign []byte,
	) (*string, []byte, error)
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
