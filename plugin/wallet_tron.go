package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
)

// tron parent
type tron struct {
	purpose  int
	coinType int

	account       uint32
	change        uint32
	addressNumber uint32

	extendedKey *keyBundle
	accountKey  *accountKey

	blockChainParams *chaincfg.Params
}

// NewAccount create new account via mnemonic wallet by
func (w *wallet) NewAccount(account, change, index uint32) (*tron, error) {
	accKey, extendedKey, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, index)
	if err != nil {
		return nil, err
	}

	return &tron{
		purpose:          Bip44Purpose,
		coinType:         TronCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    index,
		extendedKey:      extendedKey,
		accountKey:       accKey,
		blockChainParams: w.Network,
	}, err
}

// GetAddress get address with 0x
func (e *tron) GetAddress() (string, error) {
	return pubKeyToTronAddress(*e.extendedKey.PublicECDSA), nil
}

// GetPubKey get key with 0x
func (e *tron) GetPubKey() string {
	return e.extendedKey.PublicHex()
}

// GetPrvKey get key with 0x
func (e *tron) GetPrvKey() (string, error) {

	return hex.EncodeToString(e.extendedKey.PrivateECDSA.D.Bytes()), nil
}

// GetPath ...
func (e *tron) GetPath() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		e.GetPurpose(), e.GetCoinType(), e.account, e.change, e.addressNumber)
}

// GetPurpose ...
func (e *tron) GetPurpose() int {
	return e.purpose
}

// GetCoinType ...
func (e *tron) GetCoinType() int {
	return TronCoinNumber
}

func (e *tron) CloneECDSAPrivateKey() *ecdsa.PrivateKey {
	return e.extendedKey.CloneECDSAPrivateKey()
}

func (e *tron) ClearSecrets() {
	e.accountKey.ExtendedKey.Zero()
	e.accountKey.Public = ""
	e.accountKey.Private = ""

	e.extendedKey.ClearSecrets()

	e.blockChainParams = nil
	e.accountKey.Network = nil
	e.accountKey = nil
	e.extendedKey = nil
	e.account = 0
	e.change = 0
	e.addressNumber = 0
}
