package main

import (
	"github.com/tyler-smith/go-bip39"
)

func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}
