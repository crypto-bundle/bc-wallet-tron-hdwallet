package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"math/big"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
)

const addrPatKeyTemplate = "%d'/%d/%d"

type addressData struct {
	address    string
	privateKey *ecdsa.PrivateKey
}

func (e *addressData) ClonePrivateKey() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     (&big.Int{}).SetBytes(e.privateKey.X.Bytes()),
			Y:     (&big.Int{}).SetBytes(e.privateKey.Y.Bytes()),
		},
		D: (&big.Int{}).SetBytes(e.privateKey.D.Bytes()),
	}
}

type mnemonicWalletUnit struct {
	mu *sync.Mutex

	hdWalletSvc *wallet

	mnemonicWalletUUID string
	mnemonicHash       string

	// addressPool is pool of derivation addresses with private keys and address
	// map key - string with derivation path
	// map value - ecdsa.PrivateKey and address string
	addressPool map[string]*addressData
}

func (u *mnemonicWalletUnit) Shutdown(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.unloadWallet()
	if err != nil {
		return fmt.Errorf("unable to unload wallet: %w", err)
	}

	return nil
}

func (u *mnemonicWalletUnit) UnloadWallet() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.unloadWallet()
}

func (u *mnemonicWalletUnit) unloadWallet() error {
	u.hdWalletSvc.ClearSecrets()
	u.hdWalletSvc = nil

	for accountPath, data := range u.addressPool {
		if data == nil {
			continue
		}

		if data.privateKey != nil {
			zeroKey(data.privateKey)
		}

		delete(u.addressPool, accountPath)
	}

	u.addressPool = nil
	u.mnemonicWalletUUID = "0"
	u.mnemonicHash = "0"

	return nil
}

func (u *mnemonicWalletUnit) GetWalletUUID() string {
	return u.mnemonicWalletUUID
}

func (u *mnemonicWalletUnit) SignData(ctx context.Context,
	accountParameters *anypb.Any,
	dataForSign []byte,
) (*string, []byte, error) {
	accIdentity := &AccountIdentity{}
	err := proto.Unmarshal(accountParameters.GetValue(), accIdentity)
	if err != nil {
		return nil, nil, err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	return u.signData(ctx,
		accIdentity.AccountIndex,
		accIdentity.InternalIndex,
		accIdentity.AddressIndex,
		dataForSign)
}

func (u *mnemonicWalletUnit) signData(ctx context.Context,
	account, change, index uint32,
	dataForSign []byte,
) (*string, []byte, error) {
	addr, privKey, err := u.loadAddressDataByPath(ctx, account, change, index)
	if err != nil {
		return nil, nil, err
	}

	h256h := sha256.New()
	h256h.Write(dataForSign)
	hash := h256h.Sum(nil)

	signedData, err := crypto.Sign(hash, privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to sign: %w", err)
	}

	return addr, signedData, nil
}

func (u *mnemonicWalletUnit) LoadAccount(ctx context.Context,
	accountIdentity []byte,
) (*string, error) {
	accIdentity := &AccountIdentity{}
	err := proto.Unmarshal(accountIdentity, accIdentity)
	if err != nil {
		return nil, err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	addrData, _, err := u.loadAddressDataByPath(ctx, accIdentity.AccountIndex,
		accIdentity.InternalIndex,
		accIdentity.AddressIndex)
	if err != nil {
		return nil, err
	}

	if addrData == nil {
		return nil, nil
	}

	return addrData, nil
}

func (u *mnemonicWalletUnit) loadAddressDataByPath(ctx context.Context,
	account, change, index uint32,
) (*string, *ecdsa.PrivateKey, error) {
	mapKey := fmt.Sprintf(addrPatKeyTemplate, account, change, index)
	addrData, isExists := u.addressPool[mapKey]
	if !isExists {
		tronAccount, walletErr := u.hdWalletSvc.NewAccount(account, change, index)
		if walletErr != nil {
			return nil, nil, walletErr
		}

		addr, walletErr := tronAccount.GetAddress()
		if walletErr != nil {
			return nil, nil, walletErr
		}

		addrData = &addressData{
			address:    addr,
			privateKey: tronAccount.CloneECDSAPrivateKey(),
		}
		u.addressPool[mapKey] = addrData

		tronAccount.ClearSecrets()
		tronAccount = nil
	}

	return &addrData.address, addrData.ClonePrivateKey(), nil
}

func (u *mnemonicWalletUnit) GetAccountAddressByPath(ctx context.Context,
	accountIdentityRaw []byte,
) (*string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	accIdentity := &AccountIdentity{}
	err := proto.Unmarshal(accountIdentityRaw, accIdentity)
	if err != nil {
		return nil, err
	}

	return u.getAddressByPath(ctx, accIdentity.AccountIndex,
		accIdentity.InternalIndex,
		accIdentity.AddressIndex)
}

func (u *mnemonicWalletUnit) GetMultipleAccounts(ctx context.Context,
	multipleAccountsParameters *anypb.Any,
) (uint, []*pbCommon.AccountIdentity, error) {
	list := &RangeUnitsList{}
	err := proto.Unmarshal(multipleAccountsParameters.GetValue(), list)
	if err != nil {
		return 0, nil, err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	return u.getMultipleAccounts(ctx, list)
}

func (u *mnemonicWalletUnit) getMultipleAccounts(ctx context.Context,
	rangeList *RangeUnitsList,
) (uint, []*pbCommon.AccountIdentity, error) {
	var err error
	wg := sync.WaitGroup{}
	size := len(rangeList.RangeUnits)
	wg.Add(size)

	list := make([]*pbCommon.AccountIdentity, size)

	for i := uint32(0); i != uint32(size); i++ {
		rangeUnit := rangeList.RangeUnits[i]

		if rangeUnit.AddressIndexFrom == rangeUnit.AddressIndexTo { // if one item in range
			accountIdentifier, loopErr := u.getAddressAndMarshal(ctx, rangeUnit.AccountIndex,
				rangeUnit.InternalIndex, rangeUnit.AddressIndexFrom)
			if loopErr != nil {
				err = loopErr

				continue
			}

			list[i] = accountIdentifier

			wg.Done()

			continue
		}

		for addressIndex := rangeUnit.AddressIndexFrom; addressIndex <= rangeUnit.AddressIndexTo; addressIndex++ {
			go func(accountIdx, internalIdx, addressIdx, position uint32) {
				defer wg.Done()

				accountIdentifier, loopErr := u.getAddressAndMarshal(ctx, rangeUnit.AccountIndex,
					rangeUnit.InternalIndex, addressIdx)
				if loopErr != nil {
					err = loopErr

					return
				}

				list[i] = accountIdentifier

				return
			}(rangeUnit.AccountIndex, rangeUnit.InternalIndex, addressIndex, i)
		}
	}

	wg.Wait()

	if err != nil {
		return 0, nil, err
	}

	return uint(size), list, nil
}

func (u *mnemonicWalletUnit) getAddressAndMarshal(ctx context.Context,
	account, change, index uint32,
) (*pbCommon.AccountIdentity, error) {
	address, err := u.getAddressByPath(ctx, account,
		change, index)
	if err != nil {
		return nil, err
	}

	addrParams := &anypb.Any{}
	err = addrParams.MarshalFrom(&pbCommon.DerivationAddressIdentity{
		AccountIndex:  account,
		InternalIndex: change,
		AddressIndex:  index,
	})
	if err != nil {
		return nil, err
	}

	return &pbCommon.AccountIdentity{
		Parameters: addrParams,
		Address:    *address,
	}, nil
}

func (u *mnemonicWalletUnit) getAddressByPath(_ context.Context,
	account, change, index uint32,
) (*string, error) {
	tronWallet, err := u.hdWalletSvc.NewAccount(account, change, index)
	if err != nil {
		return nil, err
	}

	defer func() {
		tronWallet.ClearSecrets()
		tronWallet = nil
	}()

	blockchainAddress, err := tronWallet.GetAddress()
	if err != nil {
		return nil, err
	}

	return &blockchainAddress, nil
}

func NewPoolUnit(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error) {
	hdWalletSvc, createErr := newWalletFromMnemonic(mnemonicDecryptedData)
	if createErr != nil {
		return nil, createErr
	}

	return &mnemonicWalletUnit{
		mu: &sync.Mutex{},

		hdWalletSvc: hdWalletSvc,

		mnemonicWalletUUID: walletUUID,

		addressPool: make(map[string]*addressData),
	}, nil
}
