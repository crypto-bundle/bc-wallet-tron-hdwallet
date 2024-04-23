package grpc

import (
	"context"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type AddressForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	AccountIndex  uint32 `valid:"type(uint32)"`
	InternalIndex uint32 `valid:"type(uint32)"`
	AddressIndex  uint32 `valid:"type(uint32)"`
}

func (f *AddressForm) LoadAndValidateLoadAddrReq(ctx context.Context,
	req *pbApi.LoadDerivationAddressRequest,
) (valid bool, err error) {
	if req.MnemonicWalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}

	if req.AddressIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}

	return f.validate(req.MnemonicWalletIdentifier, req.AddressIdentifier)
}

func (f *AddressForm) LoadAndValidateGetAddrReq(ctx context.Context,
	req *pbApi.DerivationAddressRequest,
) (valid bool, err error) {
	if req.MnemonicWalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}

	if req.AddressIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Address identity")
	}

	return f.validate(req.MnemonicWalletIdentifier, req.AddressIdentity)
}

func (f *AddressForm) validate(mnemoIdentifier *pbCommon.MnemonicWalletIdentity,
	addrIdentifier *pbCommon.DerivationAddressIdentity,
) (valid bool, err error) {

	f.WalletUUID = mnemoIdentifier.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(mnemoIdentifier.WalletUUID)
	if err != nil {
		return false, err
	}

	f.AccountIndex = addrIdentifier.AccountIndex
	f.InternalIndex = addrIdentifier.InternalIndex
	f.AddressIndex = addrIdentifier.AddressIndex

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
