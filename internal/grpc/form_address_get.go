package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type GetAddressForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	AccountIndex  uint32 `valid:"type(uint)"`
	InternalIndex uint32 `valid:"type(uint)"`
	AddressIndex  uint32 `valid:"type(uint)"`
}

func (f *GetAddressForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (valid bool, err error) {
	if req.MnemonicWalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicWalletIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.MnemonicWalletIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	if req.AddressIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	f.AccountIndex = req.AddressIdentity.AccountIndex
	f.InternalIndex = req.AddressIdentity.InternalIndex
	f.AddressIndex = req.AddressIdentity.AddressIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
