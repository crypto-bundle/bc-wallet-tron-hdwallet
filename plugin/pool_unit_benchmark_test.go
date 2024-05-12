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
