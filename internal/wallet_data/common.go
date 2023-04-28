package wallet_data

import (
	"context"
	"github.com/google/uuid"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/entities"
)

type dbStoreService interface {
	AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	UpdateIsEnabledWalletByUUID(ctx context.Context, uuid string, isEnabled bool) error
	GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error)
	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.Wallet, error)
}
