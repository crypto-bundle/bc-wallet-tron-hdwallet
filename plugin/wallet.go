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

// Package hdwallet Main
package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	Bip49Purpose   = 49
	Bip44Purpose   = 44
	DefaultPurpose = Bip49Purpose

	BtcCoinNumber = 0

	TronCoinNumber = 195
	TronBytePrefix = byte(0x41)
)

// wallet contains the individual mnemonic and seed
type wallet struct {
	mnemonic string
	seed     []byte

	*keyBundle
}

// restore mnemonic a Bip32 HD-wallet for the mnemonic
func restore(mnemonic string, network *chaincfg.Params) (*wallet, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	bundle, err := newBundledKeyBySeed(seed)
	if err != nil {
		return nil, err
	}

	bundle.ExtendedKey.SetNet(network)

	return &wallet{
		mnemonic:  mnemonic,
		seed:      seed,
		keyBundle: bundle,
	}, nil
}

// Seed return seed
func (w *wallet) Seed() []byte {
	return w.seed
}

// GetSeed return string of seed from byte
func (w *wallet) GetSeed() string {
	return hex.EncodeToString(w.Seed())
}

// GetMnemonic return mnemonic string
func (w *wallet) GetMnemonic() string {
	return w.mnemonic
}

// ClearSecrets is function clear sensitive secrets data
func (w *wallet) ClearSecrets() {
	w.mnemonic = "0"

	pattern := []byte{0x1, 0x2, 0x3, 0x4}
	// Copy the pattern into the start of the container
	copy(w.seed, pattern)
	// Incrementally duplicate the pattern throughout the container
	for j := len(pattern); j < len(w.seed); j *= 2 {
		copy(w.seed[j:], w.seed[:j])
	}
	w.seed = nil

	w.keyBundle.ClearSecrets()
}

// NewWalletFromMnemonic new HD-wallet via entropy
func newWalletFromMnemonic(mnemonic string) (*wallet, error) {
	return newFromString(mnemonic, &chaincfg.MainNetParams)
}

// NewWalletFromEntropy HD-wallet via entropy
func newWalletFromEntropy(entropy []byte) (*wallet, error) {
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return restore(mnemonic, &chaincfg.MainNetParams)
}

// newFromString hdwallet via mnemo string
func newFromString(mnemo string, network *chaincfg.Params) (*wallet, error) {
	entropy, err := bip39.EntropyFromMnemonic(mnemo)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	return restore(mnemonic, network)
}
