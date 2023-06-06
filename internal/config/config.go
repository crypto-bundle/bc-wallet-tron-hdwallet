package config

import (
	"time"

	commonConfig "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-config/pkg/config"
	commonHealthcheck "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-redis/pkg/redis"
	commonVault "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

// Config for application
type Config struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonHealthcheck.HealthcheckHTTPConfig
	*VaultWrappedConfig
	*commonPostgres.PostgresConfig
	*commonNats.NatsConfig
	*commonRedis.RedisConfig
	// -------------------
	// Internal configs
	// -------------------
	*GrpcConfig
	*MnemonicConfig
	WalletManagerUnloadHotInterval      time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_HOT_INTERVAL" default:"15s"`
	WalletManagerUnloadInterval         time.Duration `envconfig:"WALLET_MANAGER_UNLOAD_INTERVAL" default:"8s"`
	WalletManagerMnemonicPerWalletCount uint8         `envconfig:"WALLET_MANAGER_MNEMONICS_PER_WALLET_COUNT" default:"3"`
}

func (c *Config) GetDefaultHotWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadHotInterval
}

func (c *Config) GetDefaultWalletUnloadInterval() time.Duration {
	return c.WalletManagerUnloadInterval
}

func (c *Config) GetMnemonicsCountPerWallet() uint8 {
	return c.WalletManagerMnemonicPerWalletCount
}

// Prepare variables to static configuration
func (c *Config) Prepare() error {
	return nil
}

func (c *Config) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
