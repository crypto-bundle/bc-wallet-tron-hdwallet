package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
)

func TestNewPoolUnit(t *testing.T) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		t.Fatalf("%s: %e", "unable to create entropy:", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("%s: %e", "unable to create mnemonic pharase from entory:", err)
	}

	_, err = NewPoolUnit(uuid.NewString(), mnemonic)
	if err != nil {
		t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", err)
	}
}

func TestMnemonicWalletUnit_GetWalletUUID(t *testing.T) {
	type testCase struct {
		WalletUUID  string
		Mnemonic    string
		AddressPath *pbCommon.DerivationAddressIdentity

		ExpectedAddress string
	}

	// WARN: DO NOT USE THESE MNEMONIC IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONIC IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONIC IN MAINNET OR TESTNET. Usage only in unit-tests
	tCase := &testCase{
		WalletUUID: uuid.NewString(),
		Mnemonic:   "seven kitten wire trap family giraffe globe access dinosaur upper forum aerobic dash segment cruise concert giant upon sniff armed rain royal firm state",
		AddressPath: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  5,
			InternalIndex: 5,
			AddressIndex:  55,
		},
		ExpectedAddress: "TRZUb6GVH922CHYty9NaFpVZWuf8GZJ3va",
	}

	poolUnitIntrf, loopErr := NewPoolUnit(tCase.WalletUUID, tCase.Mnemonic)
	if loopErr != nil {
		t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
	}

	poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
	if !ok {
		t.Fatalf("%s", "unable to cast interface to pool unit worker")
	}

	accountIdentity := &anypb.Any{}
	_ = accountIdentity.MarshalFrom(tCase.AddressPath)

	addr, loopErr := poolUnit.GetAccountAddress(context.Background(), accountIdentity)
	if loopErr != nil {
		t.Fatalf("%s: %e", "unable to get address from pool unit", loopErr)
	}

	if addr == nil {
		t.Fatalf("%s", "missing address in pool unit result")
	}

	if tCase.ExpectedAddress != *addr {
		t.Fatalf("%s", "address not equal with expected")
	}

	resultUUID := poolUnit.GetWalletUUID()
	if tCase.WalletUUID != resultUUID {
		t.Fatalf("%s", "wallet uuid not equal with expected")
	}
}

func TestMnemonicWalletUnit_GetAccountAddress(t *testing.T) {
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

	for _, tCase := range testCases {
		poolUnitIntrf, loopErr := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			t.Fatalf("%s", "unable to cast interface to pool unit worker")
		}

		accountIdentity := &anypb.Any{}
		_ = accountIdentity.MarshalFrom(tCase.AddressPath)

		addr, loopErr := poolUnit.GetAccountAddress(context.Background(), accountIdentity)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to get address from pool unit:", loopErr)
		}

		if addr == nil {
			t.Fatalf("%s", "missing address in pool unit result")
		}

		if tCase.ExpectedAddress != *addr {
			t.Fatalf("%s", "address not equal with expected")
		}
	}
}

func TestMnemonicWalletUnit_GetMultipleAccounts(t *testing.T) {
	type testCase struct {
		Mnemonic        string
		AddressPathList *pbCommon.RangeUnitsList

		ExpectedAddress map[string]*pbCommon.DerivationAddressIdentity
	}

	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	testCases := []*testCase{
		{
			Mnemonic: "web account soft juice relief green account rebel rifle gun follow thunder ski credit judge off educate round advice allow wink bitter first color",
			AddressPathList: &pbCommon.RangeUnitsList{
				RangeUnits: []*pbCommon.RangeRequestUnit{
					{AccountIndex: 3, InternalIndex: 8, AddressIndexFrom: 114, AddressIndexTo: 124},
					{AccountIndex: 155, InternalIndex: 5, AddressIndexFrom: 4, AddressIndexTo: 8},
					{AccountIndex: 2555, InternalIndex: 50, AddressIndexFrom: 250, AddressIndexTo: 255},
				},
			},
			ExpectedAddress: map[string]*pbCommon.DerivationAddressIdentity{
				"TAMcYXxBevDjYK1wVBoUsGr9F1Gc1BvgYX": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 114},
				"TAT4RTK3FKQbgevkiTATT9EnFar8hK3hzK": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 115},
				"TU88zp2gvowrvxxuskzbZhRWPcdB96yQNi": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 116},
				"THDKHgZxjEDfUWh7bQuF3eB3q2cipeqWrM": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 117},
				"TMkLVEeyLBcWPTqkxvLBpq37U7ZXBfHDhZ": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 118},
				"TSq5zADYT1GK4zedGCQT2Gh8EtWDn61iEd": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 119},
				"THtKSpFy4owX3NjNzsqobLdWda37Z5B7Gv": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 120},
				"TAAnRtjKHzhvkC1XGfDuVxkfMt6u5UdJ7P": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 121},
				"TUBsnau67pNS5yqpYSJ1SzQ1CsjKJsTPwS": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 122},
				"TKdbf5oMFYxWSLSUKnZbbdBGLg1zcypKzj": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 123},
				"TUc2Q73GuLxupQ7BCZnxwjhDiTPC6sro4A": {AccountIndex: 3, InternalIndex: 8, AddressIndex: 124},
				//
				"TE6JoSwxS1EdannPwasPPWrt6GcjjucDAU": {AccountIndex: 155, InternalIndex: 5, AddressIndex: 4},
				"TFjUm8k4BKpVAG9iqhJhihbRW2Gn454Z3r": {AccountIndex: 155, InternalIndex: 5, AddressIndex: 5},
				"TTLjtMRL2V3v2NUGbigFGejCkF5DDtLvVp": {AccountIndex: 155, InternalIndex: 5, AddressIndex: 6},
				"TGDoLJNa5oT46wSsVvFiY18YcoNXE2AaMQ": {AccountIndex: 155, InternalIndex: 5, AddressIndex: 7},
				"TZ8z3nS98FBQiZSFifLF6m9h23EDFrgA6e": {AccountIndex: 155, InternalIndex: 5, AddressIndex: 8},
				//
				"TDA3ZnCuYPJNdGiax5XfVFVh4fm3beayuk": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 250},
				"TWMgaesoun1GVbhjfKBYqvRb5PJU6eW9xk": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 251},
				"TRQdPDaeHd2ZbZGKGkASx6oX3jTPHsjkwf": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 252},
				"TCALJChrKN5Q9FjCZpYGTL5KGipE66ZF94": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 253},
				"TDpavWFBnCUbTSZF6yxk4UdGa4yPCNpaiY": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 254},
				"TXeLYYR6bE29aRA8pAfPqbu1Svx3hZgdVg": {AccountIndex: 2555, InternalIndex: 50, AddressIndex: 255},
			},
		},
	}

	for _, tCase := range testCases {
		poolUnitIntrf, loopErr := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			t.Fatalf("%s", "unable to cast interface to pool unit worker")
		}

		anyRangeUnit := &anypb.Any{}
		err := anyRangeUnit.MarshalFrom(tCase.AddressPathList)
		if err != nil {
			t.Fatalf("%s", "unable to marshal request units list")
		}

		count, addrList, loopErr := poolUnit.GetMultipleAccounts(context.Background(), anyRangeUnit)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to get address from pool unit:", loopErr)
		}

		if count == 0 {
			t.Fatalf("%s", "addr list count not equal with expected")
		}

		if addrList == nil {
			t.Fatalf("%s", "missing addr lsit in pool unit result")
		}

		if count != uint(len(addrList)) {
			t.Fatalf("%s", "length of addrlist count not equal with count rseult value")
		}

		if count != uint(len(tCase.ExpectedAddress)) {
			t.Fatalf("%s", "count value or getAccount result not equal with values count of expected map")
		}

		for i := 0; i != len(addrList); i++ {
			addr := addrList[i]

			accIdentifier, isExists := tCase.ExpectedAddress[addr.Address]
			if !isExists {
				t.Fatalf("%s", "missing addr in expected map")
			}

			marshaledAccIdentifier := &pbCommon.DerivationAddressIdentity{}
			marshalErr := addr.Parameters.UnmarshalTo(marshaledAccIdentifier)
			if marshalErr != nil {
				t.Fatalf("%s: %e", "unable to unmarshal anydata to account identity", marshalErr)
			}

			if accIdentifier.AddressIndex != marshaledAccIdentifier.AddressIndex {
				t.Fatalf("%s", "marshaled address index not equal with expected")
			}

			if accIdentifier.InternalIndex != marshaledAccIdentifier.InternalIndex {
				t.Fatalf("%s", "marshaled internal index not equal with expected")
			}

			if accIdentifier.AccountIndex != marshaledAccIdentifier.AccountIndex {
				t.Fatalf("%s", "marshaled account index not equal with expected")
			}
		}
	}
}

func TestMnemonicWalletUnit_LoadAddressByPath(t *testing.T) {
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
			Mnemonic: "umbrella uphold security hill monkey skin either immense kid afraid sense desk extend twenty doctor odor buzz reject derive frame hub much once suffer",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  5,
				InternalIndex: 12,
				AddressIndex:  3,
			},
			ExpectedAddress: "TS98RrhGNPeXNXqFYhYfeFd2AAUK7z5aED",
		},
		{
			Mnemonic: "slogan follow oil world head protect patrol wagon toddler fly kangaroo kite dash essay shoulder worth one grace shift good disease biology magic pottery",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  1000,
				InternalIndex: 10000,
				AddressIndex:  100000,
			},
			ExpectedAddress: "TZ8Tenb9okzq4x626vuux3p3SMXahu3LyG",
		},
		{
			Mnemonic: "image video differ dumb later child gather smart supply mountain salon ring boy mystery hope secret present bar then joke latin guitar view devote",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  1,
				InternalIndex: 102,
				AddressIndex:  1003,
			},
			ExpectedAddress: "TLFMBiQvLjhp9AK9N2wvgc77dsuKBuLsiV",
		},
	}

	for _, tCase := range testCases {
		poolUnitIntrf, err := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if err != nil {
			t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", err)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			t.Fatalf("%s", "unable to cast interface to pool unit worker")
		}

		accountIdentity := &anypb.Any{}
		_ = accountIdentity.MarshalFrom(tCase.AddressPath)

		addr, err := poolUnit.LoadAccount(context.Background(), accountIdentity)
		if err != nil {
			t.Fatalf("%s: %e", "unable to get address from pool unit:", err)
		}

		if addr == nil {
			t.Fatalf("%s: %e", "missing address in pool unit result:", err)
		}

		if len(poolUnit.addressPool) == 0 {
			t.Fatalf("%s", "address in pool not loaded")
		}

		key := fmt.Sprintf(addrPatKeyTemplate, tCase.AddressPath.AccountIndex,
			tCase.AddressPath.InternalIndex, tCase.AddressPath.AddressIndex)
		addrData, ok := poolUnit.addressPool[key]
		if !ok || addrData == nil {
			t.Fatalf("%s", "missing data by key in address pool")
		}

		if addrData.privateKey == nil {
			t.Fatalf("%s", "missing private key in address pool unit")
		}

		if tCase.ExpectedAddress != addrData.address {
			t.Fatalf("%s", "address not equal with expected")
		}

		if tCase.ExpectedAddress != *addr {
			t.Fatalf("%s", "address not equal with expected")
		}
	}
}

func TestMnemonicWalletUnit_SignData(t *testing.T) {
	type testCase struct {
		Mnemonic         string
		AddressPath      *pbCommon.DerivationAddressIdentity
		AddressPublicKey string
		DataForSign      []byte

		ExpectedAddress    string
		ExpectedSignedData []byte
	}

	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	testCases := []*testCase{
		{
			Mnemonic: "unknown valid carbon hat echo funny artist letter desk absorb unit fatigue foil skirt stay case path rescue hawk remember aware arch regular cry",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  7,
				InternalIndex: 8,
				AddressIndex:  9,
			},
			DataForSign: []byte{0x0, 0x2, 0x3, 0x4},

			ExpectedAddress: "TS5nGhVnjSVudb58XLUCqYTgYtjY7abJJH",
		},
		{
			Mnemonic: "laundry file mystery rate absorb wrist despair cook near afraid account mirror name chair lake regular vicious oblige release vicious identify glimpse flight help",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  909,
				InternalIndex: 8008,
				AddressIndex:  70007,
			},
			DataForSign: []byte{0x5, 0x6, 0x7, 0x8},

			ExpectedAddress: "THjwZivHc9kyosKYTF7MJLTqeCq2xdMgL6",
		},
		{
			Mnemonic: "busy spawn solar december element round wild buddy furnace help clog tired object camera resist maze fuel need stock rule spot diagram aisle expect",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  9,
				InternalIndex: 8,
				AddressIndex:  7,
			},
			AddressPublicKey: "030ba1318a2d4258cecce5725c393e3a6ab7d60cde9e6f39106cd99cf63aa36032",
			DataForSign:      []byte{0x9, 0x10, 0x11, 0x12},

			ExpectedAddress: "TPuC6Y4aFGZPu1kfiiWgkPVSUgi2oYynRz",
		},
	}

	for _, tCase := range testCases {
		poolUnitIntrf, loopErr := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			t.Fatalf("%s", "unable to cast interface to pool unit worker")
		}

		accountIdentity := &anypb.Any{}
		_ = accountIdentity.MarshalFrom(tCase.AddressPath)

		addr, signedData, loopErr := poolUnit.SignData(context.Background(), accountIdentity, tCase.DataForSign)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to sign data:", loopErr)
		}

		if addr == nil {
			t.Fatalf("%s", "missing address in result of sign method")
		}

		if signedData == nil {
			t.Fatalf("%s", "missing signed data in result of sign method")
		}

		if len(poolUnit.addressPool) == 0 {
			t.Fatalf("%s", "address in pool not loaded")
		}

		key := fmt.Sprintf(addrPatKeyTemplate, tCase.AddressPath.AccountIndex,
			tCase.AddressPath.InternalIndex, tCase.AddressPath.AddressIndex)
		addrData, ok := poolUnit.addressPool[key]
		if !ok || addrData == nil {
			t.Fatalf("%s", "missing data by key in address pool")
		}

		if addrData.privateKey == nil {
			t.Fatalf("%s", "missing private key in address pool unit")
		}

		if tCase.ExpectedAddress != addrData.address {
			t.Fatalf("%s", "address not equal with expected")
		}

		if tCase.ExpectedAddress != *addr {
			t.Fatalf("%s", "address not equal with expected")
		}

		signed := bytes.Clone(signedData)
		// DIRTY HACK
		// https://stackoverflow.com/questions/49085737/geth-ecrecover-invalid-signature-recovery-id
		// https://gist.github.com/dcb9/385631846097e1f59e3cba3b1d42f3ed#file-eth_sign_verify-go
		if signed[crypto.RecoveryIDOffset] == 27 || signed[crypto.RecoveryIDOffset] == 28 {
			signed[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
		}

		h256h := sha256.New()
		h256h.Write(tCase.DataForSign)
		hash := h256h.Sum(nil)

		h256h.Reset()
		h256h = nil

		pubKey, loopErr := crypto.SigToPub(hash, signed)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to get public key from signed message", loopErr)
		}

		tronAddr := pubKeyToTronAddress(*pubKey)
		if tCase.ExpectedAddress != tronAddr {
			t.Fatalf("%s", "tron addr from pubKey not equal with expected")
		}
	}
}

func TestMnemonicWalletUnit_UnloadWallet(t *testing.T) {
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
			Mnemonic: "input erase buzz crew miss auction habit cargo wrestle perfect like midnight buddy chase grit only treat stuff rival worth alien tennis parent artist",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  5,
				InternalIndex: 8,
				AddressIndex:  11,
			},
			ExpectedAddress: "TNnvBFnjrsdTqCnjPRsZZP24pPA1VYUqAi",
		},
		{
			Mnemonic: "empower plate axis divorce neither noodle above flight very indoor zone mango sand exhaust nominee solid combine picnic gospel myth stem raw garage veteran",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  2,
				InternalIndex: 4,
				AddressIndex:  8,
			},
			ExpectedAddress: "TWnoUuXdREoJFFc2vAuwHPUz33tuKJaonK",
		},
		{
			Mnemonic: "sea vault tattoo laugh ugly where saddle six usage install one cube affair sick used actress zebra fuel sunny tackle can siege develop drop",
			AddressPath: &pbCommon.DerivationAddressIdentity{
				AccountIndex:  8,
				InternalIndex: 64,
				AddressIndex:  4096,
			},
			ExpectedAddress: "TDXRtZoqjkJxtr68deRyKJ5Kkkf8u4kJS1",
		},
	}

	for _, tCase := range testCases {
		poolUnitIntrf, loopErr := NewPoolUnit(uuid.NewString(), tCase.Mnemonic)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to create mnemonic wallet pool unit:", loopErr)
		}

		poolUnit, ok := poolUnitIntrf.(*mnemonicWalletUnit)
		if !ok {
			t.Fatalf("%s", "unable to cast interface to pool unit worker")
		}

		accountIdentity := &anypb.Any{}
		_ = accountIdentity.MarshalFrom(tCase.AddressPath)

		addr, loopErr := poolUnit.LoadAccount(context.Background(), accountIdentity)
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to sign data:", loopErr)
		}

		if addr == nil {
			t.Fatalf("%s", "missing address in result of sign method")
		}

		if len(poolUnit.addressPool) == 0 {
			t.Fatalf("%s", "address in pool not loaded")
		}

		key := fmt.Sprintf(addrPatKeyTemplate, tCase.AddressPath.AccountIndex,
			tCase.AddressPath.InternalIndex, tCase.AddressPath.AddressIndex)
		addrData, ok := poolUnit.addressPool[key]
		if !ok || addrData == nil {
			t.Fatalf("%s", "missing data by key in address pool")
		}

		loopErr = poolUnit.UnloadWallet()
		if loopErr != nil {
			t.Fatalf("%s: %e", "unable to unload wallet", loopErr)
		}

		if len(poolUnit.addressPool) != 0 {
			t.Fatalf("%s", "address pool is not empty")
		}

		if poolUnit.hdWalletSvc != nil {
			t.Fatalf("%s", "hdwallet service is not nil")
		}

		if poolUnit.mnemonicHash != "0" {
			t.Fatalf("%s", "mnemonicHash is not equal zero")
		}
	}
}
