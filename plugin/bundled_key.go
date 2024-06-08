/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

var (
	// defaultNetwork for generate masterKey
	// nolint:gochecknoglobals // its library function
	defaultNetwork = &chaincfg.MainNetParams
	// zeroQuote base zero
	// nolint:gochecknoglobals // its library function
	zeroQuote uint32 = 0x80000000
)

// keyBundle struct
type keyBundle struct {
	filler []big.Word
	// ExtendedKey hdwallet
	ExtendedKey *hdkeychain.ExtendedKey

	// Network chain params
	Network *chaincfg.Params

	// Private for btc child's
	Private *btcec.PrivateKey
	// Public for btc child's
	Public *btcec.PublicKey

	// PrivateECDSA
	PrivateECDSA *ecdsa.PrivateKey
	// PrivateECDSA
	PublicECDSA *ecdsa.PublicKey
}

// newBundledKeyBySeed generate new extended key
func newBundledKeyBySeed(seed []byte) (*keyBundle, error) {
	extendedKey, err := hdkeychain.NewMaster(seed, defaultNetwork)
	if err != nil {
		return nil, err
	}

	bundle := &keyBundle{
		filler:      []big.Word{0, 0, 0, 0},
		ExtendedKey: extendedKey,
		Network:     defaultNetwork,
	}
	if err = bundle.init(); err != nil {
		return nil, err
	}

	return bundle, nil
}

// newBundledKeyByExtendedKey generate new bundled key
func newBundledKeyByExtendedKey(extendedKey *hdkeychain.ExtendedKey) (*keyBundle, error) {
	bundle := &keyBundle{
		filler:      []big.Word{0, 0, 0, 0},
		ExtendedKey: extendedKey,
		Network:     defaultNetwork,
	}
	if err := bundle.init(); err != nil {
		return nil, err
	}

	return bundle, nil
}

func (k *keyBundle) init() error {
	var err error

	k.Private, err = k.ExtendedKey.ECPrivKey()
	if err != nil {
		return err
	}

	k.Public, err = k.ExtendedKey.ECPubKey()
	if err != nil {
		return err
	}

	k.PrivateECDSA = k.Private.ToECDSA()
	k.PublicECDSA = &k.PrivateECDSA.PublicKey

	return nil
}

// GetPath return path in bip44 style
func (k *keyBundle) GetPath(purpose, coinType, account, change, addressIndex uint32) []uint32 {
	purpose = zeroQuote + purpose
	coinType = zeroQuote + coinType
	account = zeroQuote + account

	return []uint32{
		purpose,
		coinType,
		account,
		change,
		addressIndex,
	}
}

// GetChildKey path for address
func (k *keyBundle) GetChildKey(purpose, coinType,
	account,
	change,
	addressIndex uint32,
) (*accountKey, *keyBundle, error) {
	var err error

	extendedKeyCloned, err := k.ExtendedKey.CloneWithVersion(k.ExtendedKey.Version())
	if err != nil {
		return nil, nil, err
	}

	accKey := extendedKeyCloned
	var extendedKey = extendedKeyCloned
	for i, v := range k.GetPath(purpose, coinType, account, change, addressIndex) {
		extendedKey, err = extendedKey.Derive(v)
		if err != nil {
			return nil, nil, err
		}

		if i == 2 {
			accKey = extendedKey

			continue
		}
	}

	acc := &accountKey{
		ExtendedKey: accKey,
		Network:     k.Network,
	}

	err = acc.Init()
	if err != nil {
		return nil, nil, err
	}

	newExtendedKey, err := newBundledKeyByExtendedKey(extendedKey)
	if err != nil {
		return nil, nil, err
	}

	return acc, newExtendedKey, err
}

// PublicHex generate public key to string by hex
func (k *keyBundle) PublicHex() string {
	return hex.EncodeToString(k.Public.SerializeCompressed())
}

// PublicHash generate public key by hash160
func (k *keyBundle) PublicHash() ([]byte, error) {
	address, err := k.ExtendedKey.Address(k.Network)
	if err != nil {
		return nil, err
	}

	return address.ScriptAddress(), nil
}

// AddressP2PKH generate public key to p2wpkh style address
func (k *keyBundle) AddressP2PKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	addr1, err := btcutil.NewAddressScriptHash(script, k.Network)
	if err != nil {
		return "", err
	}

	return addr1.EncodeAddress(), nil
}

// AddressP2WPKH generate public key to p2wpkh style address
func (k *keyBundle) AddressP2WPKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

// AddressP2WPKHInP2SH generate public key to p2wpkh nested within p2sh style address
func (k *keyBundle) AddressP2WPKHInP2SH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	addr1, err := btcutil.NewAddressScriptHash(script, k.Network)
	if err != nil {
		return "", err
	}

	return addr1.EncodeAddress(), nil
}

// CloneECDSAPrivateKey full clone ECDSA private key
func (k *keyBundle) CloneECDSAPrivateKey() *ecdsa.PrivateKey {
	clonedPrivKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     (&big.Int{}).SetBytes(k.PrivateECDSA.X.Bytes()),
			Y:     (&big.Int{}).SetBytes(k.PrivateECDSA.Y.Bytes()),
		},
		D: (&big.Int{}).SetBytes(k.PrivateECDSA.D.Bytes()),
	}

	return &clonedPrivKey
}

func (k *keyBundle) ClearSecrets() {
	k.ExtendedKey.Zero()
	k.Private.Zero()

	k.PublicECDSA.X.SetBits(k.filler)
	k.PublicECDSA.X = nil
	k.PublicECDSA.Y.SetBits(k.filler)
	k.PublicECDSA.Y = nil
	k.PublicECDSA.Curve = nil

	k.PrivateECDSA.D.SetBits(k.filler)
	k.PrivateECDSA.D = nil

	k.Network = nil
	k.PrivateECDSA = nil
	k.PublicECDSA = nil
	k.Private = nil
	k.Public = nil
	k.ExtendedKey = nil

	k.filler = nil
}
