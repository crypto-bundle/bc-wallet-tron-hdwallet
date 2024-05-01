package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
)

type configService interface {
	GetHdWalletPluginPath() string
}

type walletMakerFunc func(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error)

type WalletPoolUnitService interface {
	UnloadWallet() error

	GetWalletUUID() string
	LoadAccount(ctx context.Context,
		accountIdentity []byte,
	) (*string, error)
	GetAccountAddressByPath(ctx context.Context,
		accountIdentityRaw []byte,
	) (*string, error)
	GetAddressesByPathByRange(ctx context.Context,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignData(ctx context.Context,
		accountIdentity []byte,
		dataForSign []byte,
	) (*string, []byte, error)
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}
