package config

import (
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
)

// HdWalletConfig for application
type HdWalletConfig struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonHealthcheck.HealthcheckHTTPConfig
	*pbApi.HdWalletClientConfig // yes, client config for listen on unix file socket
	*VaultWrappedConfig
	// -------------------
	// Internal configs
	// -------------------
	*MnemonicConfig
}

// Prepare variables to static configuration
func (c *HdWalletConfig) Prepare() error {
	return nil
}

func (c *HdWalletConfig) PrepareWith(cfgSvcList ...interface{}) error {
	return nil
}
