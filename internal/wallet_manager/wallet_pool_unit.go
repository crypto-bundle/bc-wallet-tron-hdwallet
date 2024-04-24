package wallet_manager

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/ethereum/go-ethereum/crypto"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/google/uuid"

	"github.com/btcsuite/btcd/chaincfg"
	"go.uber.org/zap"
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
	logger *zap.Logger

	mu *sync.Mutex

	cfgSrv       configService
	hdWalletSvc  hdWalleter
	encryptorSvc encryptService

	mnemonicEncryptedData []byte
	mnemonicWalletUUID    *uuid.UUID
	mnemonicHash          string

	// addressPool is pool of derivation addresses with private keys and address
	// map key - string with derivation path
	// map value - ecdsa.PrivateKey and address string
	addressPool map[string]*addressData
}

func (u *mnemonicWalletUnit) Init(ctx context.Context) error {
	return nil
}

func (u *mnemonicWalletUnit) Run(ctx context.Context) error {
	err := u.loadWallet(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *mnemonicWalletUnit) Shutdown(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.unloadWallet()
	if err != nil {
		u.logger.Error("unable to unload wallet", zap.Error(err))

		return err
	}

	for i := range u.mnemonicWalletUUID {
		u.mnemonicWalletUUID[i] = 0
	}

	u.mnemonicWalletUUID = nil

	u.cfgSrv = nil
	u.logger = nil

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

	for i := range u.mnemonicEncryptedData {
		u.mnemonicEncryptedData[i] = 0
	}
	u.mnemonicEncryptedData = nil

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
	account, change, index uint32,
	dataForSign []byte,
) (*string, []byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.signData(ctx, account, change, index, dataForSign)
}

func (u *mnemonicWalletUnit) signData(ctx context.Context,
	account, change, index uint32,
	dataForSign []byte,
) (address *string, signedData []byte, err error) {
	addrData, err := u.loadAddressByPath(ctx, account, change, index)
	if err != nil {
		return nil, nil, err
	}

	h256h := sha256.New()
	h256h.Write(dataForSign)
	hash := h256h.Sum(nil)

	signedData, signErr := crypto.Sign(hash, addrData.privateKey)
	if signErr != nil {
		u.logger.Error("unable to sign", zap.Error(signErr),
			zap.String(app.HDWalletAddressTag, addrData.address))

		return nil, nil, signErr
	}

	return &addrData.address, signedData, nil
}

func (u *mnemonicWalletUnit) LoadAddressByPath(ctx context.Context,
	account, change, index uint32,
) (*string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	addrData, err := u.loadAddressByPath(ctx, account, change, index)
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
		tronWallet, walletErr := u.hdWalletSvc.NewTronWallet(account, change, index)
		if walletErr != nil {
			return nil, walletErr
		}

		addr, walletErr := tronWallet.GetAddress()
		if walletErr != nil {
			return nil, walletErr
		}

		clonedPrivKey, walletErr := tronWallet.ExtendedKey.CloneECDSAPrivateKey()
		if walletErr != nil {
			return nil, walletErr
		}

		addrData = &addressData{
			address:    addr,
			privateKey: clonedPrivKey,
		}

		u.addressPool[key] = addrData

		// clear temporary keys
		// TODO: add Clear method to hdwallet.Tron instance - tronWallet
		//defer func() {
		//	zeroKey(tronWallet.ExtendedKey.PrivateECDSA)
		//	zeroPubKey(tronWallet.ExtendedKey.PublicECDSA)
		//}()
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
				u.logger.Error("unable to get address by path", zap.Error(getAddrErr),
					zap.Uint32(app.HDWalletAccountIndexTag, rangeUnit.AccountIndex),
					zap.Uint32(app.HDWalletInternalIndexTag, rangeUnit.InternalIndex),
					zap.Uint32(app.HDWalletAddressIndexTag, rangeUnit.InternalIndex))

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
					u.logger.Error("unable to get address by path", zap.Error(getAddrErr),
						zap.Uint32(app.HDWalletAccountIndexTag, rangeUnit.AccountIndex),
						zap.Uint32(app.HDWalletInternalIndexTag, rangeUnit.InternalIndex),
						zap.Uint32(app.HDWalletAddressIndexTag, addressIdx))

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
	tronWallet, err := u.hdWalletSvc.NewTronWallet(account, change, index)
	if err != nil {
		return nil, err
	}

	blockchainAddress, err := tronWallet.GetAddress()
	if err != nil {
		return nil, err
	}

	return &blockchainAddress, nil
}

func (u *mnemonicWalletUnit) loadWallet(ctx context.Context) error {
	mnemonicBytes, err := u.encryptorSvc.Decrypt(u.mnemonicEncryptedData)
	if err != nil {
		return err
	}

	mnemonicSum256 := sha256.Sum256(mnemonicBytes)
	u.mnemonicHash = hex.EncodeToString(mnemonicSum256[:])

	blockChainParams := chaincfg.MainNetParams
	hdWallet, creatErr := hdwallet.NewFromString(string(mnemonicBytes), &blockChainParams)
	if creatErr != nil {
		return creatErr
	}
	u.hdWalletSvc = hdWallet

	for i := range u.mnemonicEncryptedData {
		u.mnemonicEncryptedData[i] = 0
	}
	u.mnemonicEncryptedData = nil

	return nil
}

func newMnemonicWalletPoolUnit(logger *zap.Logger,
	walletUUID uuid.UUID,
	encryptorSvc encryptService,
	mnemonicEncryptedData []byte,
) *mnemonicWalletUnit {
	return &mnemonicWalletUnit{
		logger: logger.With(zap.String(app.MnemonicWalletUUIDTag, walletUUID.String())),
		mu:     &sync.Mutex{},

		hdWalletSvc: nil, // that field will be field @ load wallet stage

		encryptorSvc: encryptorSvc,

		mnemonicWalletUUID:    &walletUUID,
		mnemonicEncryptedData: mnemonicEncryptedData,

		addressPool: make(map[string]*addressData),
	}
}
