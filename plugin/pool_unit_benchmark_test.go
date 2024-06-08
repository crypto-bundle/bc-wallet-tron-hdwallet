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
	"context"
	"testing"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

var allocsPool = make(map[string]*mnemonicWalletUnit)

func BenchmarkMnemonicWalletUnit_GetAccountAddress(b *testing.B) {
	type testCase struct {
		Mnemonic    string
		AddressPath *pbCommon.DerivationAddressIdentity

		ExpectedAddress string
	}

	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	testCases := []*testCase{
		{
			Mnemonic: "unfair silver dune air rib enforce protect limit jazz dinner thumb drift spring warrior bonus snack argue flavor wild faculty derive open dynamic carpet",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  3,
				InternalIndex: 13,
				AddressIndex:  114,
			},
			ExpectedAddress: "TFyMUdJsREv3Q1ooMhV5r2UDGFSL4xgFeC",
		},
		{
			Mnemonic: "obscure town quick bundle north message want sketch brass tone vast spoil home gentle field ozone mushroom current math cat canvas plunge stay truly",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  1020,
				InternalIndex: 10300,
				AddressIndex:  104000,
			},
			ExpectedAddress: "TBKmbAG6JefDEg741YpsMPTB7MegySqs45",
		},
		{
			Mnemonic: "beach large spray gentle buyer hover flock dream hybrid match whip ten mountain pitch enemy lobster afford barrel patrol desk trigger output excuse truck",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  2,
				InternalIndex: 104,
				AddressIndex:  1005,
			},
			ExpectedAddress: "TWW7CQdsogbqfc5FrSb6MKu22QS4Reg3mH",
		},
	}

	b.ReportAllocs()

	for _, tCase := range testCases {
		poolUnitIntrf, loopErr := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if loopErr != nil {
			b.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			b.Fatalf("%s", "unable to cast interface to pool unit worker")
		}
		allocsPool[poolUnit.GetWalletUUID()] = poolUnit

		accountIdentity := &anypb.Any{}
		_ = accountIdentity.MarshalFrom(tCase.AddressPath)

		addr, loopErr := poolUnit.GetAccountAddress(context.Background(), accountIdentity)
		if loopErr != nil {
			b.Fatalf("%s: %e", "unable to get address from pool unit:", loopErr)
		}

		if addr == nil {
			b.Fatalf("%s", "missing address in pool unit result")
		}

		if tCase.ExpectedAddress != *addr {
			b.Fatalf("%s", "address not equal with expected")
		}

		loopErr = poolUnit.Shutdown(context.Background())
		if loopErr != nil {
			b.Fatalf("%s", "unable to shurdown pool unit")
		}
	}
}
