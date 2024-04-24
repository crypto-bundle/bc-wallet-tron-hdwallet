package grpc

import (
	"context"
	"fmt"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type UnLoadMnemonicForm struct {
	WalletUUID    string    `valid:"type(string),uuid,required"`
	WalletUUIDRaw uuid.UUID `valid:"-"`
}

func (f *UnLoadMnemonicForm) LoadAndValidate(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (valid bool, err error) {
	if req.MnemonicIdentity == nil {
		return false, fmt.Errorf("%w:%s", ErrMissedRequiredData, "Wallet identity")
	}
	f.WalletUUID = req.MnemonicIdentity.WalletUUID
	f.WalletUUIDRaw, err = uuid.Parse(req.MnemonicIdentity.WalletUUID)
	if err != nil {
		return false, err
	}

	_, err = govalidator.ValidateStruct(f)
	if err != nil {
		return false, err
	}

	return true, nil
}
