package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type SignDataForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	AccountIndex  uint32 `valid:"type(uint32)"`
	InternalIndex uint32 `valid:"type(uint32)"`
	AddressIndex  uint32 `valid:"type(uint32)"`

	DataForSign []byte `valid:"required"`
}

func (f *SignDataForm) LoadAndValidate(ctx context.Context,
	req *pbApi.SignDataRequest,
) (valid bool, err error) {
	if req.MnemonicWalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicWalletIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.MnemonicWalletIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	if req.AddressIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}
	f.AccountIndex = req.AddressIdentifier.AccountIndex
	f.InternalIndex = req.AddressIdentifier.InternalIndex
	f.AddressIndex = req.AddressIdentifier.AddressIndex

	f.DataForSign = req.DataForSign

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
