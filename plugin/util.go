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

// Package hdwallet implements heirarchical deterministic Bitcoin wallets, as defined in BIP 32.
//
// BIP 32 - https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
//
// This package provides utilities for generating hierarchical deterministic Bitcoin wallets.
//
// Examples
//
//          // Generate a random 256 bit seed
//          seed, _ := genSeed(256)
//
//          // Create a master private key
//          masterPrv := masterKey(seed)
//
//          // Convert a private key to public key
//          masterPub := masterPrv.Pub()
//
//          // Generate new child key based on private or public key
//          childPrv, err := masterPrv.Child(0)
//          childPub, err := masterPrv.Child(0)
//
//          // Create bitcoin address from public key
//          address := childPub.Address()
//
//          // Convenience string -> string Child and Address functions
//          walletString := childPub.String()
//          childString, _ := stringChild(walletString,0)
//          childAddress, _ := stringAddress(childString)
//
// Extended Keys
//
// Hierarchical deterministic wallets are simply deserialized extended keys. Extended Keys can be imported
// and exported as base58-encoded strings. Here are two examples:

// public key:
// "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"

// private key:
// "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"
package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shengdoushi/base58"
	// nolint:staticcheck // its library function
	"golang.org/x/crypto/ripemd160"
)

// nolint:gochecknoglobals // its library function
var curve *btcec.KoblitzCurve = btcec.S256()

// nolint:deadcode,unused // its library function
func stringToMagic(magic string) [4]byte {
	var b [4]byte
	t, _ := hex.DecodeString(magic)
	copy(b[:], t)

	return b
}

func pubKeyToTronAddress(key ecdsa.PublicKey) string {
	addr := crypto.PubkeyToAddress(key)

	addrTrxBytes := make([]byte, 0)
	addrTrxBytes = append(addrTrxBytes, TronBytePrefix)
	addrTrxBytes = append(addrTrxBytes, addr.Bytes()...)

	crc := calcCheckSum(addrTrxBytes)

	addrTrxBytes = append(addrTrxBytes, crc...)

	//nolint:gocritic // ok. Just reminder hot to generate address_hex for TronKey table
	// addrTrxHex := hex.EncodeToString(addrTrxBytes)

	addrTrx := base58.Encode(addrTrxBytes, base58.BitcoinAlphabet)

	return addrTrx
}

func hash160(data []byte) ([]byte, error) {
	sha := sha256.New()
	ripe := ripemd160.New()
	_, err := sha.Write(data)
	if err != nil {
		return nil, err
	}

	_, err = ripe.Write(sha.Sum(nil))
	if err != nil {
		return nil, err
	}

	return ripe.Sum(nil), nil
}

func dblSha256(data []byte) ([]byte, error) {
	sha1 := sha256.New()
	sha2 := sha256.New()
	_, err := sha1.Write(data)
	if err != nil {
		return nil, err
	}

	_, err = sha2.Write(sha1.Sum(nil))
	if err != nil {
		return nil, err
	}

	return sha2.Sum(nil), nil
}

func privToPub(key []byte) []byte {
	return compress(curve.ScalarBaseMult(key))
}

func onCurve(x, y *big.Int) bool {
	return curve.IsOnCurve(x, y)
}

func compress(x, y *big.Int) []byte {
	two := big.NewInt(2)
	rem := two.Mod(y, two).Uint64()
	rem += 2
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(rem))
	rest := x.Bytes()
	pad := 32 - len(rest)
	if pad != 0 {
		zeroes := make([]byte, pad)
		rest = append(zeroes, rest...)
	}
	return append(b[1:], rest...)
}

// 2.3.4 of SEC1 - http://www.secg.org/index.php?action=secg,docs_secg
// nolint: gocritic // cuz of library function
func expand(key []byte) (*big.Int, *big.Int) {
	params := curve.Params()
	exp := big.NewInt(1)
	exp.Add(params.P, exp)
	exp.Div(exp, big.NewInt(4))
	x := big.NewInt(0).SetBytes(key[1:33])
	y := big.NewInt(0).SetBytes(key[:1])
	beta := big.NewInt(0)
	beta.Exp(x, big.NewInt(3), nil)
	beta.Add(beta, big.NewInt(7))
	beta.Exp(beta, exp, params.P)
	if y.Add(beta, y).Mod(y, big.NewInt(2)).Int64() == 0 {
		y = beta
	} else {
		y = beta.Sub(params.P, beta)
	}

	return x, y
}

func addPrivKeys(k1, k2 []byte) []byte {
	i1 := big.NewInt(0).SetBytes(k1)
	i2 := big.NewInt(0).SetBytes(k2)
	i1.Add(i1, i2)
	i1.Mod(i1, curve.Params().N)
	k := i1.Bytes()
	zero, _ := hex.DecodeString("00")
	return append(zero, k...)
}

func addPubKeys(k1, k2 []byte) []byte {
	x1, y1 := expand(k1)
	x2, y2 := expand(k2)
	return compress(curve.Add(x1, y1, x2, y2))
}

func uint32ToByte(i uint32) []byte {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)
	return a
}

func uint16ToByte(i uint16) []byte {
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, i)
	return a[1:]
}

func byteToUint16(b []byte) uint16 {
	if len(b) == 1 {
		zero := make([]byte, 1)
		b = append(zero, b...)
	}
	return binary.BigEndian.Uint16(b)
}

func calcCheckSum(data []byte) []byte {
	h256h0 := sha256.New()
	h256h0.Write(data)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	return h1[:4]
}

func zeroKey(key *ecdsa.PrivateKey) {
	tpl := []big.Word{0x0}

	key.D.SetBits(tpl)
	key.X.SetBits(tpl)
	key.Y.SetBits(tpl)
}
