package plugin

import (
	"fmt"
	"plugin"
)

func stringFuncSymbolLookUp(plugin *plugin.Plugin, symbolName string) (func() string, error) {
	pluginFuncSymbol, err := plugin.Lookup(getPluginNameSymbol)
	if err != nil {
		return nil, err
	}

	pluginFunc, ok := pluginFuncSymbol.(func() string)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnableCastPluginEntry, symbolName)
	}

	return pluginFunc, nil
}
