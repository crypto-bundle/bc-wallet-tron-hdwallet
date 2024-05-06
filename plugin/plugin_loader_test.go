package main

import (
	"fmt"
	"plugin"
	"strconv"
	"testing"
	"time"
)

func TestPlugin_Loader(t *testing.T) {
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

	p, err := plugin.Open("../build/tron.so")
	if err != nil {
		t.Fatalf("%s: %e", "unable to load pluggable extension", err)
	}

	getPluginReleaseTagFunc, err := stringFuncSymbolLookUp(p, getPluginNameSymbol)
	if err != nil {
		t.Fatal(err)
	}

	if getPluginReleaseTagFunc == nil {
		t.Fatal("missing Get release tag function")
	}

	releaseTag := getPluginReleaseTagFunc()
	if len(releaseTag) == 0 {
		t.Fatal("zero length of release tag value")
	}

	getPluginCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginCommitIDSymbol)
	if err != nil {
		t.Fatal(err)
	}

	if getPluginCommitIDFunc == nil {
		t.Fatal("missing Get commit id function")
	}

	commitID := getPluginReleaseTagFunc()
	if len(commitID) != 40 {
		t.Fatal("wrong length of commit ID value")
	}

	getPluginShortCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginShortCommitIDSymbol)
	if err != nil {
		t.Fatal(err)
	}

	if getPluginShortCommitIDFunc == nil {
		t.Fatal("missing Get short commit id function")
	}

	shortCommitID := getPluginReleaseTagFunc()
	if len(shortCommitID) != 12 {
		t.Fatal("wrong length of short commit ID value")
	}

	getPluginBuildNumberFunc, err := stringFuncSymbolLookUp(p, getPluginBuildNumberSymbol)
	if err != nil {
		t.Fatal(err)
	}

	if getPluginBuildNumberFunc == nil {
		t.Fatal("missing Get plugin build number function")
	}

	buildNumber := getPluginBuildNumberFunc()
	if _, err = strconv.ParseUint(buildNumber, 10, 0); err != nil {
		t.Fatalf("%s: %e", "wrong build number format", err)
	}

	getPluginBuildDateTSFunc, err := stringFuncSymbolLookUp(p, getPluginBuildDateTSSymbol)
	if err != nil {
		t.Fatal(err)
	}

	if getPluginBuildDateTSFunc == nil {
		t.Fatal("missing Get plugin build date function")
	}

	buildDate := getPluginBuildDateTSFunc()
	buildDataTS, err := strconv.ParseUint(buildDate, 10, 0)
	if err != nil {
		t.Fatalf("%s: %e", "wrong build date Time stamp format", err)
	}
	tm := time.Unix(int64(buildDataTS), 0)
	tmString := strconv.FormatUint(uint64(tm.Unix()), 0)
	if buildDate != tmString {
		t.Fatalf("%s, current value: %s", "something wrong with build date time stamp string",
			buildDate)
	}
}
