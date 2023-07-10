# Change Log

## [v0.0.24] - 28.04.2023 17:50 MSK

### Changed

#### Switching to a proprietary license.
License of **bc-wallet-tron-hdwallet** repository changed to proprietary - commit revision number **323a3a909a30e3ff455672da1a7bd630e936c9e4**.

Origin repository - https://github.com/crypto-bundle/bc-wallet-tron-hdwallet

The MIT license is replaced by me (_Kotelnikov Aleksei_) as an author and maintainer.

The license has been replaced with a proprietary one, with the condition of maintaining the authorship
and specifying in the README.md file in the section of authors and contributors.

[@gudron (Kotelnikov Aleksei)](https://github.com/gudron) - author and maintainer of [crypto-bundle project](https://github.com/crypto-bundle)

The commit is signed with the key -
gudron2s@gmail.com
E456BB23A18A9347E952DBC6655133DD561BF3EC

## [v0.0.26] - 14.05.2023

### Changed
* Docker container build
* Fixed Nats-kv cache bucket bugs

## [v0.0.27] - 09.06.2023
### Added 
* Logs in main.go

## [v0.0.28] - 13.06.2023
### Fixed
* Create wallet bug. Bug in restoration from cache MnemonicWalletItem entity

## [v0.0.29] - 06.07.2023
### Fixed
* Application can't process init stage with empty wallets table

## [v0.0.30] - 11.07.2023
### Added
* sync.Pool usage in GetDerivationAddress gRPC method
* AddressIdentitiesCount parameter to GetDerivationAddressByRange response
### Changed
* GetDerivationAddressByRange gRPC method - added support of multiple ranges per request
* Added iterator pattern to request form - DerivationAddressByRangeForm