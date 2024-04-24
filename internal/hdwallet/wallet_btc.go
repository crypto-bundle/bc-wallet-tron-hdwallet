package hdwallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// NewBtcWallet create new wallet
func (w *Wallet) NewBtcWallet(account, change, address uint32) (*BTC, error) {
	blockChainParams := chaincfg.MainNetParams

	blockChainParams.HDPrivateKeyID = [4]byte{0x04, 0x9d, 0x78, 0x78} // yprv
	blockChainParams.HDPublicKeyID = [4]byte{0x04, 0x9d, 0x7c, 0xb2}  // ypub

	w.Network = &blockChainParams
	accountKey, key, err := w.GetChildKey(DefaultPurpose, BtcCoinNumber, account, change, address)
	if err != nil {
		return nil, err
	}

	return &BTC{
		blockChainParams: &blockChainParams,
		purpose:          DefaultPurpose,
		coinType:         BtcCoinNumber,
		account:          account,
		change:           change,
		addressNumber:    address,

		ExtendedKey: key,
		AccountKey:  accountKey,
	}, nil
}

// GetAddress get address string
func (b *BTC) GetAddress() (string, error) {
	return b.ExtendedKey.AddressP2WPKHInP2SH()
}

// GetP2WPKHAddress get address string
func (b *BTC) GetP2WPKHAddress() (string, error) {
	return b.ExtendedKey.AddressP2WPKH()
}

// GetP2SHAddress get address string
func (b *BTC) GetP2SHAddress() (string, error) {
	return b.ExtendedKey.AddressP2WPKHInP2SH()
}

// GetP2PKHAddress get address string
func (b *BTC) GetP2PKHAddress() (string, error) {
	return b.ExtendedKey.AddressP2PKH()
}

// GetPrvKey get address private key
func (b *BTC) GetPrvKey() (string, error) {
	prvKey, err := btcutil.NewWIF(b.ExtendedKey.Private, b.ExtendedKey.Network, true)
	if err != nil {
		return "", nil
	}
	return prvKey.String(), nil
}

// GetWIF get address private key
func (b *BTC) GetWIF() (*btcutil.WIF, error) {
	prvKey, err := btcutil.NewWIF(b.ExtendedKey.Private, b.ExtendedKey.Network, true)
	if err != nil {
		return nil, err
	}

	return prvKey, nil
}

// GetPubKey get address public key
func (b *BTC) GetPubKey() string {
	return b.ExtendedKey.PublicHex()
}

// AccountPrvKey return string key
func (b *BTC) AccountPrvKey() string {
	return b.AccountKey.Private
}

// AccountPrvKeyNoMagic return string key
func (b *BTC) AccountPrvKeyNoMagic() string {
	return b.AccountKey.Private[4:]
}

// AccountPubKey return string key
func (b *BTC) AccountPubKey() string {
	return b.AccountKey.Public
}

// AccountPubKeyNoMagic return string key
func (b *BTC) AccountPubKeyNoMagic() string {
	return b.AccountKey.Public[4:]
}

// PKH return string key
func (b *BTC) PKH() (string, error) {
	return b.ExtendedKey.AddressP2WPKH()
}

// HEX return string key
func (b *BTC) HEX() string {
	return b.ExtendedKey.PublicHex()
}

// HASH return string key
func (b *BTC) HASH() ([]byte, error) {
	return b.ExtendedKey.PublicHash()
}

// GetAccountWIF get address private key
func (b *BTC) GetAccountWIF() (*btcutil.WIF, error) {
	ecKey, err := b.AccountKey.ExtendedKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.NewWIF(ecKey, b.ExtendedKey.Network, true)
	if err != nil {
		return nil, err
	}

	return wif, nil
}

// GetPurpose
func (b *BTC) GetPurpose() int {
	return b.purpose
}

// GetBlockChainConfig
func (b *BTC) GetBlockChainConfig() *chaincfg.Params {
	return b.blockChainParams
}

// GetPath
func (b *BTC) GetPath() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		b.GetPurpose(), b.GetCoinType(), b.account, b.change, b.addressNumber)
}

// GetCoinNumber
func (b *BTC) GetCoinType() int {
	return b.coinType
}
