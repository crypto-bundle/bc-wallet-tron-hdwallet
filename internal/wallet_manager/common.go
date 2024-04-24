package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"
)

type configService interface {
}

type WalletPoolUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	GetMnemonicUUID() *uuid.UUID
	LoadAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (*string, error)
	GetAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (*string, error)
	GetAddressesByPathByRange(ctx context.Context,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		account, change, index uint32,
		dataForSign []byte,
	) (*string, []byte, error)
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewTronWallet(account, change, address uint32) (*hdwallet.Tron, error)

	ClearSecrets()
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}
