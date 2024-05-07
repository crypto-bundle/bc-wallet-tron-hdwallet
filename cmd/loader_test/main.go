package main

import (
	"context"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"os"
	"plugin"
	"strconv"
	"time"
)

func main() {
	type walletPoolUnitService interface {
		UnloadWallet() error

		GetWalletUUID() string
		LoadAccount(ctx context.Context,
			accountParameters *anypb.Any,
		) (*string, error)
		GetAccountAddress(ctx context.Context,
			accountParameters *anypb.Any,
		) (*string, error)
		GetMultipleAccounts(ctx context.Context,
			multipleAccountsParameters *anypb.Any,
		) (uint, []*pbCommon.AccountIdentity, error)
		SignData(ctx context.Context,
			accountParameters *anypb.Any,
			dataForSign []byte,
		) (*string, []byte, error)
	}

	const (
		getPluginNameSymbol          = "GetPluginName"
		getPluginReleaseTagSymbol    = "GetPluginReleaseTag"
		getPluginCommitIDSymbol      = "GetPluginCommitID"
		getPluginShortCommitIDSymbol = "GetPluginShortCommitID"
		getPluginBuildNumberSymbol   = "GetPluginBuildNumber"
		getPluginBuildDateTSSymbol   = "GetPluginBuildDateTS"

		pluginGenerateMnemonicSymbol = "GenerateMnemonic"
		pluginValidateMnemonicSymbol = "ValidateMnemonic"
		pluginNewPoolUnitSymbol      = "NewPoolUnit"
	)

	stringFuncSymbolLookUp := func(plugin *plugin.Plugin, symbolName string) (func() string, error) {
		pluginFuncSymbol, err := plugin.Lookup(symbolName)
		if err != nil {
			return nil, err
		}

		pluginFunc, ok := pluginFuncSymbol.(func() string)
		if !ok {
			return nil, fmt.Errorf("%s: %s", "unable to cast plugin symbol", symbolName)
		}

		return pluginFunc, nil
	}

	p, err := plugin.Open("./build/tron.so")
	if err != nil {
		log.Fatalf("%s: %e", "unable to load pluggable extension", err)
	}

	log.Printf("=== RUN: %s\n", getPluginNameSymbol)

	getPluginReleaseTagFunc, err := stringFuncSymbolLookUp(p, getPluginReleaseTagSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if getPluginReleaseTagFunc == nil {
		log.Fatal("missing Get release tag function")
	}

	releaseTag := getPluginReleaseTagFunc()
	if len(releaseTag) == 0 {
		log.Fatalf("%s: %d", "zero length of release tag value", len(releaseTag))
	}

	log.Printf("--- PASS: %s\n", getPluginNameSymbol)

	log.Printf("=== RUN: %s\n", getPluginCommitIDSymbol)

	getPluginCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginCommitIDSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if getPluginCommitIDFunc == nil {
		log.Fatalf("missing Get commit id function")
	}

	commitID := getPluginCommitIDFunc()
	if len(commitID) != 40 {
		log.Fatalf("%s: %d, %s", "wrong length of commit ID value", len(commitID), commitID)
	}

	log.Printf("--- PASS: %s\n", getPluginCommitIDSymbol)

	log.Printf("=== RUN: %s\n", getPluginShortCommitIDSymbol)

	getPluginShortCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginShortCommitIDSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if getPluginShortCommitIDFunc == nil {
		log.Fatal("missing Get short commit id function")
	}

	shortCommitID := getPluginShortCommitIDFunc()
	if len(shortCommitID) != 7 {
		log.Fatalf("%s: %d, %s", "wrong length of short commit ID value", len(shortCommitID),
			shortCommitID)
	}

	log.Printf("--- PASS: %s\n", getPluginShortCommitIDSymbol)

	log.Printf("=== RUN: %s\n", getPluginBuildNumberSymbol)

	getPluginBuildNumberFunc, err := stringFuncSymbolLookUp(p, getPluginBuildNumberSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if getPluginBuildNumberFunc == nil {
		log.Fatal("missing Get plugin build number function")
	}

	buildNumber := getPluginBuildNumberFunc()
	if _, err = strconv.ParseUint(buildNumber, 10, 0); err != nil {
		log.Fatalf("%s: %e", "wrong build number format", err)
	}

	log.Printf("--- PASS: %s\n", getPluginBuildNumberSymbol)

	log.Printf("=== RUN: %s\n", getPluginBuildDateTSSymbol)

	getPluginBuildDateTSFunc, err := stringFuncSymbolLookUp(p, getPluginBuildDateTSSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if getPluginBuildDateTSFunc == nil {
		log.Fatal("missing Get plugin build date function")
	}

	buildDate := getPluginBuildDateTSFunc()
	buildDataTS, err := strconv.ParseUint(buildDate, 10, 0)
	if err != nil {
		log.Fatalf("%s: %e", "wrong build date Time stamp format", err)
	}
	tm := time.Unix(int64(buildDataTS), 0)
	tmString := strconv.FormatUint(uint64(tm.Unix()), 10)
	if buildDate != tmString {
		log.Fatalf("%s, current value: %s", "something wrong with build date time stamp string",
			buildDate)
	}

	log.Printf("--- PASS: %s\n", getPluginBuildDateTSSymbol)

	log.Printf("=== RUN: %s\n", pluginGenerateMnemonicSymbol)

	generateMnemonicFuncSymbol, err := p.Lookup(pluginGenerateMnemonicSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if generateMnemonicFuncSymbol == nil {
		log.Fatal("missing Generate mnemonic function symbol")
	}

	generateMnemoFunc, ok := generateMnemonicFuncSymbol.(func() (string, error))
	if !ok {
		log.Fatal("unable to cast generate mnemonic function")
	}

	generatedMnemonic, err := generateMnemoFunc()
	if err != nil {
		log.Fatal(err)
	}

	if !bip39.IsMnemonicValid(generatedMnemonic) {
		log.Fatal("generated mnemonic phares is not valid")
	}

	log.Printf("--- PASS: %s\n", pluginGenerateMnemonicSymbol)

	log.Printf("=== RUN: %s\n", pluginValidateMnemonicSymbol)

	validateMnemonicFuncSymbol, err := p.Lookup(pluginValidateMnemonicSymbol)
	if err != nil {
		log.Fatal(err)
	}

	if validateMnemonicFuncSymbol == nil {
		log.Fatal("missing Validate mnemonic function symbol")
	}

	validateMnemoFunc, isCasted := validateMnemonicFuncSymbol.(func(mnemonic string) bool)
	if !isCasted {
		log.Fatal("unable to cast validate function")
	}

	if validateMnemoFunc == nil {
		log.Fatal("missing Validate mnemonic function")
	}

	// WARN: DO NOT USE THESE MNEMONICS IN MAINNET OR TESTNET. Usage only in unit-tests
	mnemoPhrase := "beach large spray gentle buyer hover flock dream hybrid match whip ten mountain pitch enemy lobster afford barrel patrol desk trigger output excuse truck"
	if !validateMnemoFunc(mnemoPhrase) {
		log.Fatal("failed mnemonic validation validation")
	}

	log.Printf("--- PASS: %s\n", pluginValidateMnemonicSymbol)

	log.Printf("=== RUN: %s\n", pluginNewPoolUnitSymbol)

	unitMakerFuncSymbol, err := p.Lookup(pluginNewPoolUnitSymbol)
	if err != nil {
		log.Fatal(err)
	}

	unitMakerFunc, isCasted := unitMakerFuncSymbol.(func(walletUUID string,
		mnemonicDecryptedData string,
	) (interface{}, error))
	if !isCasted {
		log.Fatal("unable to cast pool unit Maker function")
	}

	unitInterface, err := unitMakerFunc(uuid.NewString(), mnemoPhrase)
	if err != nil {
		log.Fatal(err)
	}

	_, isCasted = unitInterface.(walletPoolUnitService)
	if !isCasted {
		log.Fatal("unable to cast pool unit Maker to named interface")
	}

	log.Printf("--- PASS: %s\n", pluginNewPoolUnitSymbol)

	log.Println("PASS")

	os.Exit(0)
}
