package forms

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

type DerivationAddressByRangeForm struct {
	WalletUUID            string `valid:"type(string),uuid,required"`
	WalletUUIDRaw         uuid.UUID
	MnemonicWalletUUID    string `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID

	AccountIndex     uint32 `valid:"type(uint32),int,required"`
	InternalIndex    uint32 `valid:"type(uint32),int,required"`
	AddressIndexFrom uint32 `valid:"type(uint32),int,required"`
	AddressIndexTo   uint32 `valid:"type(uint32),int,required"`
}

func (f *DerivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	if req.WalletIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.WalletIdentity.WalletUUID

	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "MnemonicWallet identity")
	}
	f.MnemonicWalletUUID = req.MnemonicIdentity.WalletUUID

	f.AccountIndex = req.AccountIndex
	f.InternalIndex = req.InternalIndex
	f.AddressIndexFrom = req.AddressIndexFrom
	f.AddressIndexTo = req.AddressIndexTo

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
