package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"github.com/google/uuid"
)

type addressRangeForm struct {
	AccountIndex     uint32 `valid:"type(uint)"`
	InternalIndex    uint32 `valid:"type(uint)"`
	AddressIndexFrom uint32 `valid:"type(uint)"`
	AddressIndexTo   uint32 `valid:"type(uint)"`
	AddressIndexDiff int32  `valid:"type(int32)"`
}

type derivationAddressByRangeForm struct {
	MnemonicWalletUUID    string    `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID `valid:"-"`

	Ranges      []*addressRangeForm `valid:"required"`
	RangesCount uint32              `valid:"type(uint32),required"`
	RangeSize   uint32              `valid:"type(uint32),required"`

	index uint32
}

func (f *derivationAddressByRangeForm) hasNext() bool {
	if f.index < f.RangesCount {
		return true
	}
	return false
}

func (f *derivationAddressByRangeForm) GetRangesCount() uint32 {
	return f.RangesCount
}

func (f *derivationAddressByRangeForm) GetRangesSize() uint32 {
	return f.RangeSize
}

func (f *derivationAddressByRangeForm) GetNext() *addressRangeForm {
	if f.hasNext() {
		rageForm := f.Ranges[f.index]
		f.index++

		return rageForm
	}
	return nil
}

func (f *derivationAddressByRangeForm) LoadAndValidate(ctx context.Context,
	req *pbApi.DerivationAddressByRangeRequest,
) (valid bool, err error) {
	if req.MnemonicWalletIdentifier == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.MnemonicWalletUUID = req.MnemonicWalletIdentifier.WalletUUID

	if req.Ranges == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Ranges data")
	}
	f.RangesCount = uint32(len(req.Ranges))
	f.Ranges = make([]*addressRangeForm, len(req.Ranges))
	for i := uint32(0); i != f.RangesCount; i++ {
		data := req.Ranges[i]
		diff := (data.AddressIndexTo - data.AddressIndexFrom) + 1
		if data.AddressIndexTo == data.AddressIndexFrom {
			diff = 1
		}

		f.Ranges[i] = &addressRangeForm{
			AccountIndex:     data.AccountIndex,
			InternalIndex:    data.InternalIndex,
			AddressIndexFrom: data.AddressIndexFrom,
			AddressIndexTo:   data.AddressIndexTo,
			AddressIndexDiff: int32(diff),
		}
		f.RangeSize += diff
	}

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	mnemonicWalletUUIDRaw, err := uuid.Parse(f.MnemonicWalletUUID)
	if err != nil {
		return false, err
	}
	f.MnemonicWalletUUIDRaw = mnemonicWalletUUIDRaw

	return true, nil
}
