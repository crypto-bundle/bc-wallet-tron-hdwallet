package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shengdoushi/base58"
)

// Tron parent
type Tron struct {
	purpose  int
	coinType int

	account       uint32
	change        uint32
	addressNumber uint32

	ExtendedKey *Key
	AccountKey  *accountKey

	blockChainParams *chaincfg.Params
}

// NewAccount create new account via mnemonic wallet by
func (w *Wallet) NewAccount(account, change, index uint32) (*Tron, error) {
	accKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, index)
	if err != nil {
		return nil, err
	}

	return &Tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    index,
		ExtendedKey:      key,
		AccountKey:       accKey,
		blockChainParams: w.Network,
	}, err
}

func (w *Wallet) GetAccountAddress(account, change, address uint32) (string, error) {
	accountKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, address)
	if err != nil {
		return "nil", err
	}

	tronAccountData := &Tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    address,
		ExtendedKey:      key,
		AccountKey:       accountKey,
		blockChainParams: w.Network,
	}

	return tronAccountData.GetAddress()
}

func (w *Wallet) GetAccountPublicKey(account, change, address uint32) (string, error) {
	accountKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, address)
	if err != nil {
		return "nil", err
	}

	tronAccountData := &Tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    address,
		ExtendedKey:      key,
		AccountKey:       accountKey,
		blockChainParams: w.Network,
	}

	return tronAccountData.GetPubKey(), nil
}

func (w *Wallet) GetAccountPrivateKey(account, change, address uint32) (string, error) {
	accountKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, address)
	if err != nil {
		return "nil", err
	}

	tronAccountData := &Tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    address,
		ExtendedKey:      key,
		AccountKey:       accountKey,
		blockChainParams: w.Network,
	}

	return tronAccountData.GetPrvKey()
}

func (w *Wallet) SignByPrivateKey(privateKey []byte, dataForSign []byte) ([]byte, error) {
	privKey, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(dataForSign)
	hash := h256h.Sum(nil)

	signedData, err := crypto.Sign(hash, privKey)
	if err != nil {
		return nil, err
	}

	return signedData, nil
}

func (w *Wallet) SignByAccount(account, change, address uint32, dataForSign []byte) ([]byte, error) {
	accountKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, address)
	if err != nil {
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(dataForSign)
	hash := h256h.Sum(nil)

	tronAccountData := &Tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    address,
		ExtendedKey:      key,
		AccountKey:       accountKey,
		blockChainParams: w.Network,
	}

	signedData, err := crypto.Sign(hash, tronAccountData.ExtendedKey.PrivateECDSA)
	if err != nil {
		return nil, err
	}

	return signedData, nil
}

// GetAddress get address with 0x
func (e *Tron) GetAddress() (string, error) {
	addr := crypto.PubkeyToAddress(*e.ExtendedKey.PublicECDSA)

	addrTrxBytes := make([]byte, 0)
	addrTrxBytes = append(addrTrxBytes, TronBytePrefix)
	addrTrxBytes = append(addrTrxBytes, addr.Bytes()...)

	crc := calcCheckSum(addrTrxBytes)

	addrTrxBytes = append(addrTrxBytes, crc...)

	//nolint:gocritic // ok. Just reminder hot to generate address_hex for TronKey table
	// addrTrxHex := hex.EncodeToString(addrTrxBytes)

	addrTrx := base58.Encode(addrTrxBytes, base58.BitcoinAlphabet)

	return addrTrx, nil
}

// GetPubKey get key with 0x
func (e *Tron) GetPubKey() string {
	return e.ExtendedKey.PublicHex()
}

// GetPrvKey get key with 0x
func (e *Tron) GetPrvKey() (string, error) {
	return hex.EncodeToString(e.ExtendedKey.PrivateECDSA.D.Bytes()), nil
}

// GetPath ...
func (e *Tron) GetPath() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		e.GetPurpose(), e.GetCoinType(), e.account, e.change, e.addressNumber)
}

// GetPurpose ...
func (e *Tron) GetPurpose() int {
	return e.purpose
}

// GetCoinType ...
func (e *Tron) GetCoinType() int {
	return TronCoinNumber
}

func (e *Tron) ClearSecrets() {
	e.AccountKey.ExtendedKey.Zero()
	e.AccountKey.Public = ""
	e.AccountKey.Private = ""

	e.ExtendedKey.ClearSecrets()

	e.blockChainParams = nil
	e.AccountKey.Network = nil
	e.ExtendedKey = nil
	e.account = 0
	e.change = 0
	e.addressNumber = 0
}
