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

type MnemonicWalletUnit struct {
	logger *zap.Logger

	mu sync.Mutex

	cfgSrv       configService
	hdWalletSrv  hdWalleter
	encryptorSvc encryptService

	mnemonicEncryptedData []byte
	mnemonicWalletUUID    uuid.UUID
	mnemonicHash          string

	// addressPool is pool of derivation addresses with private keys and address
	// map key - string with derivation path
	// map value - ecdsa.PrivateKey and address string
	addressPool map[string]*addressData
}

func (u *MnemonicWalletUnit) Init(ctx context.Context) error {
	return nil
}

func (u *MnemonicWalletUnit) Run(ctx context.Context) error {
	err := u.loadWallet(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *MnemonicWalletUnit) Shutdown(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.unloadWallet()
	if err != nil {
		u.logger.Error("unable to unload wallet", zap.Error(err))

		return err
	}

	return nil
}

func (u *MnemonicWalletUnit) UnloadWallet() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.unloadWallet()
}

func (u *MnemonicWalletUnit) unloadWallet() error {
	u.hdWalletSrv.ClearSecrets()
	u.hdWalletSrv = nil

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

	for i := range u.mnemonicWalletUUID {
		u.mnemonicWalletUUID[i] = 0
	}

	u.logger.Info("wallet successfully unload")
	u.logger = nil
	u.mnemonicHash = ""

	return nil
}

func (u *MnemonicWalletUnit) GetMnemonicUUID() uuid.UUID {
	return u.mnemonicWalletUUID
}

func (u *MnemonicWalletUnit) GetWalletIdentity() *pbCommon.MnemonicWalletIdentity {
	return &pbCommon.MnemonicWalletIdentity{
		WalletUUID: u.mnemonicWalletUUID.String(),
		WalletHash: u.mnemonicHash,
	}
}

func (u *MnemonicWalletUnit) SignData(ctx context.Context,
	account, change, index uint32,
	transactionData []byte,
) (*string, []byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.signData(ctx, account, change, index, transactionData)
}

func (u *MnemonicWalletUnit) signData(ctx context.Context,
	account, change, index uint32,
	transactionData []byte,
) (address *string, signedData []byte, err error) {
	addrData, err := u.loadAddressByPath(ctx, account, change, index)
	if err != nil {
		return nil, nil, err
	}

	h256h := sha256.New()
	h256h.Write(transactionData)
	hash := h256h.Sum(nil)

	signedData, signErr := crypto.Sign(hash, addrData.privateKey)
	if signErr != nil {
		u.logger.Error("unable to sign", zap.Error(signErr),
			zap.String(app.HDWalletAddressTag, addrData.address))

		return nil, nil, signErr
	}

	return &addrData.address, signedData, nil
}

func (u *MnemonicWalletUnit) LoadAddressByPath(ctx context.Context,
	account, change, index uint32,
) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	addrData, err := u.loadAddressByPath(ctx, account, change, index)
	if err != nil {
		return "", err
	}

	return addrData.address, nil
}

func (u *MnemonicWalletUnit) loadAddressByPath(ctx context.Context,
	account, change, index uint32,
) (*addressData, error) {
	key := fmt.Sprintf("%d'/%d/%d", account, change, index)
	addrData, isExists := u.addressPool[key]
	if !isExists {
		tronWallet, walletErr := u.hdWalletSrv.NewTronWallet(account, change, index)
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

func (u *MnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	account, change, index uint32,
) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.loadWallet(ctx)
	if err != nil {
		return "", err
	}

	return u.getAddressByPath(ctx, account, change, index)
}

func (u *MnemonicWalletUnit) GetAddressesByPathByRange(ctx context.Context,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.getAddressesByPathByRange(ctx, rangeIterable, marshallerCallback)
}

func (u *MnemonicWalletUnit) getAddressesByPathByRange(ctx context.Context,
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
				position, address)

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

				marshallerCallback(accountIdx, internalIdx, addressIdx, position, address)

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

func (u *MnemonicWalletUnit) getAddressByPath(_ context.Context,
	account, change, index uint32,
) (string, error) {
	tronWallet, err := u.hdWalletSrv.NewTronWallet(account, change, index)
	if err != nil {
		return "", err
	}

	blockchainAddress, err := tronWallet.GetAddress()
	if err != nil {
		return "", err
	}

	return blockchainAddress, nil
}

func (u *MnemonicWalletUnit) LoadWallet(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.loadWallet(ctx)
}

func (u *MnemonicWalletUnit) loadWallet(ctx context.Context) error {
	mnemonicBytes, err := u.encryptorSvc.Decrypt(u.mnemonicEncryptedData)
	if err != nil {
		return err
	}

	mnemonicSum256 := sha256.Sum256(mnemonicBytes)
	if hex.EncodeToString(mnemonicSum256[:]) != u.mnemonicHash {
		return ErrWrongMnemonicHash
	}

	blockChainParams := chaincfg.MainNetParams
	hdWallet, creatErr := hdwallet.NewFromString(string(mnemonicBytes), &blockChainParams)
	if creatErr != nil {
		return creatErr
	}
	u.hdWalletSrv = hdWallet

	u.logger.Info("wallet successfully load")

	return nil
}

func newMnemonicWalletPoolUnit(logger *zap.Logger,
	walletUUID uuid.UUID,
	walletHash string,
	encryptorSvc encryptService,
) *MnemonicWalletUnit {
	return &MnemonicWalletUnit{
		logger: logger.With(zap.String(app.MnemonicWalletUUIDTag, walletUUID.String())),
		mu:     sync.Mutex{},

		hdWalletSrv: nil, // that field will be field @ load wallet stage

		encryptorSvc: encryptorSvc,

		mnemonicWalletUUID: walletUUID,
		mnemonicHash:       walletHash,

		addressPool: make(map[string]*addressData),
	}
}
