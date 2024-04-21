package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type UnLoadMultipleMnemonicForm struct {
	MnemonicsList        []*UnLoadMnemonicForm `valid:"required"`
	MnemonicWalletsUUIDs []uuid.UUID           `valid:"-"`

	WalletsCount uint `valid:"-"`
}

func (f *UnLoadMultipleMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.UnLoadMultipleMnemonicsRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identities")
	}

	f.MnemonicsList = make([]*UnLoadMnemonicForm, len(req.MnemonicIdentity))
	f.MnemonicWalletsUUIDs = make([]uuid.UUID, len(req.MnemonicIdentity))

	for i, _ := range req.MnemonicIdentity {
		walletIdentityForm := &UnLoadMnemonicForm{
			WalletUUID: req.MnemonicIdentity[i].WalletUUID,
		}

		_, err = govalidator.ValidateStruct(walletIdentityForm)
		if err != nil {
			return false, err
		}

		walletIdentityForm.WalletUUIDRaw, _ = uuid.Parse(walletIdentityForm.WalletUUID)
		f.MnemonicsList[i] = walletIdentityForm
		f.MnemonicWalletsUUIDs[i] = walletIdentityForm.WalletUUIDRaw
	}

	f.WalletsCount = uint(len(req.MnemonicIdentity))

	return true, nil
}
