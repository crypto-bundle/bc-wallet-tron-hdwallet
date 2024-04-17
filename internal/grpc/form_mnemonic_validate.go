package grpc

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	"github.com/google/uuid"
)

type ValidateMnemonicForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`

	EncryptedMnemonicData []byte `valid:"required"`
}

func (f *ValidateMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.ValidateMnemonicRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.MnemonicIdentity.WalletUUID)
	if err != nil {
		return false, err
	}

	f.EncryptedMnemonicData = req.MnemonicData

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
