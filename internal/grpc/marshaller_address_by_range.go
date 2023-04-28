package grpc

import (
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/types"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallGetAddressByRange(
	walletPublicData *types.PublicWalletData,
	mnemonicWalletPublicData *types.PublicMnemonicWalletData,
	addressesData []*pbApi.DerivationAddressIdentity,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	response := &pbApi.DerivationAddressByRangeResponse{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: walletPublicData.UUID.String(),
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletPublicData.UUID.String(),
			WalletHash: mnemonicWalletPublicData.Hash,
		},
		AddressIdentities: addressesData,
	}

	return response, nil
}
