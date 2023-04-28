package grpc

import (
	"context"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/types"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

type walletManagerService interface {
	CreateNewWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title string,
		purpose string,
	) (*types.PublicWalletData, error)
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*types.PublicDerivationAddressData, error)

	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		accountIndex uint32,
		internalIndex uint32,
		addressIndexFrom uint32,
		addressIndexTo uint32,
		marshallerCallback func(addressIdx, position uint32, address string),
	) error

	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error)
	GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error)

	SignTransaction(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transaction *tronCore.Transaction,
	) (*types.PublicSignTxData, error)
}

type marshallerService interface {
	MarshallCreateWalletData(*types.PublicWalletData) (*pbApi.AddNewWalletResponse, error)
	MarshallGetAddressData(
		walletPublicData *types.PublicWalletData,
		mnemonicWalletPublicData *types.PublicMnemonicWalletData,
		addressPublicData *types.PublicDerivationAddressData,
	) (*pbApi.DerivationAddressResponse, error)
	MarshallGetAddressByRange(
		walletPublicData *types.PublicWalletData,
		mnemonicWalletPublicData *types.PublicMnemonicWalletData,
		addressesData []*pbApi.DerivationAddressIdentity,
	) (*pbApi.DerivationAddressByRangeResponse, error)
	MarshallGetEnabledWallets([]*types.PublicWalletData) (*pbApi.GetEnabledWalletsResponse, error)
	MarshallSignTransaction(
		publicSignTxData *types.PublicSignTxData,
	) (*pbApi.SignTransactionResponse, error)
	MarshallWalletInfo(
		walletData *types.PublicWalletData,
	) *pbApi.WalletData
}
