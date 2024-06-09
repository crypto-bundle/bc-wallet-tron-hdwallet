# Change Log

## [v0.0.31] 09.06.2024
### Added
* Added new plugin functions:
  * GetCoinID
  * GetSupportedCoinIDs
* Added plugin init flow - for set pluginChainID variable
  * Value of pluginChainID variable depends on build-time variable 

## [v0.0.30] 08.06.2024
### Changed
* Changed MIT License to NON-AI MIT
* Changed version of hdwallet-controller dependency to v0.0.29

## [v0.0.29] 08.06.2024
### Changed
* New version of bc-wallet-common-hdwallet-controller dependency - v0.0.28
* Changed Helm chart for:
  * Support of ManagerAPI nad WalletAPI of hdwallet-controller application
  * Support of PoW-shield

## [v0.0.28] 12.05.2024
### Added
* Added usage of lib-profiler
* Added lib-profiler config ENV ot Helm-chart
* Added golang ENV variable - GOMAXPROCS and GOMEMLIMIT
* Added usage of bc-wallet-tron-hdwallet-migration container to Helm-chart
### Changed
* Added support of new version bc-wallet-common-hdwallet-controller v0.0.27
* Changed unix-socket file ENV of gRPC listen flow
* Changed build of tron-hdwallet-api docker container - now golang:1.22.2-bookworm is stage 1 container
### Fixed
* Fixed Redis host and port ENV names

## [v0.0.27] 07.05.2024
### Changed
* Bc-wallet-tron-hdwallet-api moved to another repository
  * Now it's a common application for all hd-wallet support blockchain - bc-wallet-common-hdwallet-api
* Fixed bugs in Helm-chart
* Changed Docker container build flow:
  * Changed stages in docker file
* Changed content README.md file

## [v0.0.26] 05.05.2024
### Added
* Added plugin path ENV variable
* Added ldflags support in build application flow
* Added unit-tests to all mnemonicWalletUnit methods
* Added plugin wrapper package
* Added ldflags support in build plugin flow
* Added support of new version bc-wallet-common-hdwallet-controller v0.0.25
### Changed
* Refactored hd-wallet service-component
  * All struct and variables is un-exportable
  * Moved to tron plugin directory
* Changed gRPC methods:
  * GetDerivationAddress replaced by GetAccount
  * GetDerivationAddressByRange replaced by GetMultipleAccounts
* Refactored wallet pool unit service-component
  * Changed code for support plugin flow
  * Moved to tron plugin directory
  * Now plugin must support next exported functions:
    * NewPoolUnit
    * GenerateMnemonic
    * ValidateMnemonic
    * GetPluginName
    * GetPluginReleaseTag
    * GetPluginCommitID
    * GetPluginShortCommitID
    * GetPluginBuildNumber
    * GetPluginBuildDateTS

## [v0.0.25] 24.04.2024
### Added
* Implemented new gRPC methods:
  * GenerateMnemonic
  * LoadMnemonic
  * UnLoadMnemonic
  * EncryptMnemonic
  * ValidateMnemonic
  * UnLoadMultipleMnemonics
  * LoadDerivationAddress
  * SignData
### Changed
* Removed all types of store - postgres, redis, nats. Hdwallet-api app now storage less application
* Removed multiple mnemonics per wallet flow - now one wallet - one mnemonic
* Changed mnemonic management flow:
  * Now all wallet only in memory
  * Logic of load/unload wallet in gRPC method requests
  * Added vault encrypt/decrypt flow for passing mnemonic data between hdwallet-controller and hdwallet-api
* Bump go version 1.19 -> 1.22
* Integrated new version of hdwallet-controller dependency - v0.0.24
* Bump common-lib version:
  * bc-wallet-common-lib-config v0.0.5
  * bc-wallet-common-lib-grpc v0.0.4
  * bc-wallet-common-lib-healthcheck v0.0.4
  * bc-wallet-common-lib-logger v0.0.4
  * bc-wallet-common-lib-nats-queue v0.1.12
  * bc-wallet-common-lib-postgres v0.0.8
  * bc-wallet-common-lib-redis v0.0.7
  * bc-wallet-common-lib-tracer v0.0.4
  * bc-wallet-common-lib-vault v0.0.13

## [v0.0.24] 14.02.2024
### Info
Start of big application refactoring
### Added
* Added wallet sessions entities for storing in persistent and cache stores
### Changed
* Separated application on two parts
  * bc-wallet-common-hdwallet-controller
  * bc-wallet-tron-hdwallet
* Changed GetDerivationAddressByRange gRPC method - now support get addresses by multiple ranges
* Added HdWallet API proto description
  * new gRPC method - GenerateMnemonic
  * new gRPC method - LoadMnemonic
  * new gRPC method - UnLoadMnemonic
* Added Controller API proto description
  * new gRPC method - StartWalletSession
  * new gRPC method - GetWalletSession
* Removed go-tron-sdk dependency

## [v0.0.23 - v0.0.6] 05.04.2023 - 28.04.2023
### Added
* Added gRPC client wrapper
* Small security improvements:
  * Filling private keys with zeroes - private key clearing
* Added data cache flow for storing wallet in redis and nats
* Added new gRPC-handler - GetWalletInfo
### Changed
* Changed deployment flow
  * Added helm-chart option for docker container repository
  * Fixed helm-chart template for VAULT_DATA_PATH variable
* Optimization in get addresses by range flow
* Clone private key in sign transaction flow
### Fixed
* Fixed bug in sign transaction flow
* Fixed migrations - wrong rollback SQL-code, missing drop index and drop table

## [v0.0.5] 05.04.2023
### Added
* Encryption private data via hashicorp vault
* Added gRPC client config
### Changed
* Cleaned up repository:
  * Removed ansible database deployment script
  * Removed vault polices
  * Removed private data from helm-chart
* Updated common-libs:
  * removed old bc-wallet-common dependency
  * integrated lib-common dependencies:
    * lib-postgres
    * lib-config
    * lib-grpc
    * lib-tracer
    * lib-logger
    * lib-vault
### Fixed
* Fixed bug in wallet init stage
* Fixed crash in wallet pool init stage
* Fixed bugs in flow in new wallet creation

## [v0.0.4] 16.03.2023
### Changed
* Refactoring wallet pool service-component:
  * Added wallet pool unit
  * Added unit maker
  * Added support of multiple and single mnemonic wallet
  * Added timer for mnemonic unloading flow

## [initial] - 06.03.2022 - 16.03.2023
* Created go module as bc-wallet-eth-hdwallet
* Added proto files for gRPC API
* Integrated common dependencies
* Moved to crypto-bundle namespace
* Added wallet persistent store
* Added functionality for gRPC handlers
  * AddNewWallet
  * GetDerivationAddress
  * GetEnabledWallets
  * GetDerivationAddressByRange
* Added mnemonic encryption via rsa-keys
* Added MIT licence
* Refactoring service for supporting TRON blockchain
* Created Helm chart