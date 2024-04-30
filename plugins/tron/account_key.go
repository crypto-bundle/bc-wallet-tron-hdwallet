package main

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

// accountKey results info account keys
type accountKey struct {
	Network     *chaincfg.Params
	ExtendedKey *hdkeychain.ExtendedKey
	Private     string
	Public      string
}

// Init account keys
func (a *accountKey) Init() error {
	a.Private = a.ExtendedKey.String()
	w, _ := hdWalletFromString(a.Private, a.Network.HDPrivateKeyID, a.Network.HDPublicKeyID)
	pub, err := w.Pub().String()
	if err != nil {
		return err
	}

	a.Public = pub

	return nil
}
