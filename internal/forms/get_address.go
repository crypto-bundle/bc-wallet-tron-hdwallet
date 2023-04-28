package forms

import (
	"context"
	"errors"
	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

var (
	ErrUnableReadGrpcMetadata          = errors.New("unable to read grpc metadata")
	ErrUnableGetWalletUUIDFromMetadata = errors.New("unable to get wallet uuid from metadata")
)

type GetDerivationAddressForm struct {
	WalletUUID            string `valid:"type(string),uuid,required"`
	WalletUUIDRaw         uuid.UUID
	MnemonicWalletUUID    string `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID

	AccountIndex  uint32 `valid:"type(uint32),int,required"`
	InternalIndex uint32 `valid:"type(uint32),int,required"`
	AddressIndex  uint32 `valid:"type(uint32),int,required"`
}

func (f *GetDerivationAddressForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (valid bool, err error) {
	f.WalletUUID = req.WalletIdentity.WalletUUID
	f.MnemonicWalletUUID = req.MnemonicIdentity.WalletUUID

	f.AccountIndex = req.AddressIdentity.AccountIndex
	f.InternalIndex = req.AddressIdentity.InternalIndex
	f.AddressIndex = req.AddressIdentity.AddressIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	walletUUIDRaw, err := uuid.Parse(f.WalletUUID)
	if err != nil {
		return false, err
	}
	f.WalletUUIDRaw = walletUUIDRaw

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
