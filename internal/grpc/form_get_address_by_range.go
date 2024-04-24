package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"
)

type derivationAddressByRangeForm struct {
	MnemonicWalletUUID    string    `valid:"type(string),uuid,required"`
	MnemonicWalletUUIDRaw uuid.UUID `valid:"-"`

	Ranges      []*types.PublicDerivationAddressRangeData `valid:"required"`
	RangesCount uint32                                    `valid:"type(uint32),required"`
	RangeSize   uint32                                    `valid:"type(uint32),required"`

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

func (f *derivationAddressByRangeForm) GetNext() *types.PublicDerivationAddressRangeData {
	if f.hasNext() {
		rageForm := f.Ranges[f.index]
		f.index++

		return &types.PublicDerivationAddressRangeData{
			AccountIndex:     rageForm.AccountIndex,
			InternalIndex:    rageForm.InternalIndex,
			AddressIndexFrom: rageForm.AddressIndexFrom,
			AddressIndexTo:   rageForm.AddressIndexTo,
			AddressIndexDiff: int32(rageForm.AddressIndexDiff),
		}
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
	f.Ranges = make([]*types.PublicDerivationAddressRangeData, len(req.Ranges))
	for i := uint32(0); i != f.RangesCount; i++ {
		data := req.Ranges[i]
		var diff = int32(data.AddressIndexTo-data.AddressIndexFrom) + 1
		if data.AddressIndexTo == data.AddressIndexFrom {
			diff = 1
		}

		f.Ranges[i] = &types.PublicDerivationAddressRangeData{
			AccountIndex:     data.AccountIndex,
			InternalIndex:    data.InternalIndex,
			AddressIndexFrom: data.AddressIndexFrom,
			AddressIndexTo:   data.AddressIndexTo,
			AddressIndexDiff: diff,
		}
		f.RangeSize += uint32(diff)
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
