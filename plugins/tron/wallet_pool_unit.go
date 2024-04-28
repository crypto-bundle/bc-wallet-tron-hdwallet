package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"google.golang.org/protobuf/proto"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
)

const (
	derivationPathTemplate = "%d'/%d/%d"
)

type addressData struct {
	address    string
	privateKey *ecdsa.PrivateKey
}

type mnemonicWalletUnit struct {
	mu *sync.Mutex

	hdWalletSvc *Wallet

	mnemonicWalletUUID *uuid.UUID
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
		return fmt.Errorf("unable to unload wallet: %w")
	}

	for i := range u.mnemonicWalletUUID {
		u.mnemonicWalletUUID[i] = 0
	}

	u.mnemonicWalletUUID = nil

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

	for key, data := range u.addressPool {
		if data == nil {
			continue
		}

		if data.privateKey != nil {
			zeroKey(data.privateKey)
		}

		delete(u.addressPool, key)
	}

	u.addressPool = nil

	return nil
}

func (u *mnemonicWalletUnit) GetMnemonicUUID() *uuid.UUID {
	cloned := uuid.UUID(u.mnemonicWalletUUID[:])

	return &cloned
}

func (u *mnemonicWalletUnit) GetWalletIdentity() *pbCommon.MnemonicWalletIdentity {
	return &pbCommon.MnemonicWalletIdentity{
		WalletUUID: u.mnemonicWalletUUID.String(),
		WalletHash: u.mnemonicHash,
	}
}

func (u *mnemonicWalletUnit) SignData(ctx context.Context,
	accountIdentities []byte,
	dataForSign []byte,
) (*string, []byte, error) {
	accData := &pbCommon.DerivationAddressIdentity{}
	err := proto.Unmarshal(accountIdentities, accData)
	if err != nil {
		return nil, nil, err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	return u.signData(ctx,
		accData.AccountIndex, accData.InternalIndex, accData.AccountIndex,
		dataForSign)
}

func (u *mnemonicWalletUnit) signData(ctx context.Context,
	account, change, index uint32,
	dataForSign []byte,
) (*string, []byte, error) {
	addrData, err := u.loadAddressByPath(ctx, account, change, index)
	if err != nil {
		return nil, nil, err
	}

	h256h := sha256.New()
	h256h.Write(dataForSign)
	hash := h256h.Sum(nil)

	signedData, err := crypto.Sign(hash, addrData.privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to sign: %w", err)
	}

	return &addrData.address, signedData, nil
}

func (u *mnemonicWalletUnit) LoadAddressByPath(ctx context.Context,
	accountIdentities []byte,
) (*string, error) {
	accData := &pbCommon.DerivationAddressIdentity{}
	err := proto.Unmarshal(accountIdentities, accData)
	if err != nil {
		return nil, err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	addrData, err := u.loadAddressByPath(ctx, accData.AccountIndex, accData.InternalIndex, accData.AccountIndex)
	if err != nil {
		return nil, err
	}

	if addrData == nil {
		return nil, nil
	}

	return &addrData.address, nil
}

func (u *mnemonicWalletUnit) loadAddressByPath(ctx context.Context,
	account, change, index uint32,
) (*addressData, error) {
	key := fmt.Sprintf("%d'/%d/%d", account, change, index)
	addrData, isExists := u.addressPool[key]
	if !isExists {
		tronAccount, walletErr := u.hdWalletSvc.NewAccount(account, change, index)
		if walletErr != nil {
			return nil, walletErr
		}

		addr, walletErr := tronAccount.GetAddress()
		if walletErr != nil {
			return nil, walletErr
		}

		clonedPrivKey, walletErr := tronAccount.ExtendedKey.CloneECDSAPrivateKey()
		if walletErr != nil {
			return nil, walletErr
		}

		addrData = &addressData{
			address:    addr,
			privateKey: clonedPrivKey,
		}

		u.addressPool[key] = addrData

		tronAccount.ClearSecrets()
		tronAccount = nil
	}

	return addrData, nil
}

func (u *mnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	account, change, index uint32,
) (*string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.getAddressByPath(ctx, account, change, index)
}

func (u *mnemonicWalletUnit) GetAddressesByPathByRange(ctx context.Context,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.getAddressesByPathByRange(ctx, rangeIterable, marshallerCallback)
}

func (u *mnemonicWalletUnit) getAddressesByPathByRange(ctx context.Context,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(int(rangeIterable.GetRangesSize()))

	position := uint32(0)
	for {
		rangeUnit := rangeIterable.GetNext()
		if rangeUnit == nil {
			break
		}

		if rangeUnit.AddressIndexFrom == rangeUnit.AddressIndexTo { // if one item in range
			address, getAddrErr := u.getAddressByPath(ctx, rangeUnit.AccountIndex,
				rangeUnit.InternalIndex, rangeUnit.AddressIndexFrom)
			if getAddrErr != nil {
				err = getAddrErr

				continue
			}

			marshallerCallback(rangeUnit.AccountIndex, rangeUnit.InternalIndex, rangeUnit.AddressIndexFrom,
				position, *address)

			wg.Done()

			continue
		}

		for addressIndex := rangeUnit.AddressIndexFrom; addressIndex <= rangeUnit.AddressIndexTo; addressIndex++ {
			go func(accountIdx, internalIdx, addressIdx, position uint32) {
				defer wg.Done()

				address, getAddrErr := u.getAddressByPath(ctx, rangeUnit.AccountIndex,
					rangeUnit.InternalIndex, addressIdx)
				if getAddrErr != nil {
					err = getAddrErr
					return
				}

				marshallerCallback(accountIdx, internalIdx, addressIdx, position, *address)

				return
			}(rangeUnit.AccountIndex, rangeUnit.InternalIndex, addressIndex, position)

			position++
		}
	}

	wg.Wait()

	if err != nil {
		return err
	}

	return nil
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

func NewPoolUnit(walletUUID uuid.UUID,
	mnemonicDecryptedData []byte,
) (*mnemonicWalletUnit, error) {
	blockChainParams := chaincfg.MainNetParams
	hdWalletSvc, createErr := NewFromString(string(mnemonicDecryptedData), &blockChainParams)
	if createErr != nil {
		return nil, createErr
	}

	return &mnemonicWalletUnit{
		mu: &sync.Mutex{},

		hdWalletSvc: hdWalletSvc,

		mnemonicWalletUUID: &walletUUID,

		addressPool: make(map[string]*addressData),
	}, nil
}
