package plugin

import (
	"context"
	"fmt"
	"plugin"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
)

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

type wrapper struct {
	pluginPath string
	pluginName string

	generateFunc   generateMnemonicFunc
	validateFunc   validateMnemonicFunc
	walletMakerClb walletMakerFunc

	ldFlagManager
}

func (w *wrapper) GetPluginName() string {
	return w.pluginName
}

func (w *wrapper) GetMnemonicGeneratorFunc() func() (string, error) {
	return w.generateFunc
}

func (w *wrapper) GetMnemonicValidatorFunc() func(mnemonic string) bool {
	return w.validateFunc
}

func (w *wrapper) GetMakeWalletCallback() func(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error) {
	return w.walletMakerClb
}

func (w *wrapper) Init(_ context.Context) error {
	p, err := plugin.Open(w.pluginPath)
	if err != nil {
		return err
	}

	getPluginNameFunc, err := stringFuncSymbolLookUp(p, getPluginNameSymbol)
	if err != nil {
		return err
	}

	getPluginReleaseTagFunc, err := stringFuncSymbolLookUp(p, getPluginReleaseTagSymbol)
	if err != nil {
		return err
	}

	getPluginCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginCommitIDSymbol)
	if err != nil {
		return err
	}

	getPluginShortCommitIDFunc, err := stringFuncSymbolLookUp(p, getPluginShortCommitIDSymbol)
	if err != nil {
		return err
	}

	getPluginBuildNumberFunc, err := stringFuncSymbolLookUp(p, getPluginBuildNumberSymbol)
	if err != nil {
		return err
	}

	getPluginBuildDateTSFunc, err := stringFuncSymbolLookUp(p, getPluginBuildDateTSSymbol)
	if err != nil {
		return err
	}

	generateMnemonicFuncSymbol, err := p.Lookup(pluginGenerateMnemonicSymbol)
	if err != nil {
		return err
	}

	generateMnemoFunc, ok := generateMnemonicFuncSymbol.(func() (string, error))
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginGenerateMnemonicSymbol)
	}

	validateMnemonicFuncSymbol, err := p.Lookup(pluginValidateMnemonicSymbol)
	if err != nil {
		return err
	}

	validateMnemoFunc, ok := validateMnemonicFuncSymbol.(func(mnemonic string) bool)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginValidateMnemonicSymbol)
	}

	unitMakerFuncSymbol, err := p.Lookup(pluginNewPoolUnitSymbol)
	if err != nil {
		return err
	}

	unitMakerFunc, ok := unitMakerFuncSymbol.(func(walletUUID string,
		mnemonicDecryptedData string,
	) (interface{}, error))
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, pluginNewPoolUnitSymbol)
	}

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(getPluginReleaseTagFunc(),
		getPluginCommitIDFunc(), getPluginShortCommitIDFunc(),
		getPluginBuildNumberFunc(), getPluginBuildDateTSFunc())
	if err != nil {
		return err
	}

	w.generateFunc = generateMnemoFunc
	w.validateFunc = validateMnemoFunc
	w.ldFlagManager = flagManagerSvc
	w.pluginName = getPluginNameFunc()
	w.walletMakerClb = unitMakerFunc

	return nil
}

func NewPlugin(pluginPath string) *wrapper {
	return &wrapper{
		pluginPath:     pluginPath,
		pluginName:     "",
		walletMakerClb: nil,
		ldFlagManager:  nil,
	}
}
