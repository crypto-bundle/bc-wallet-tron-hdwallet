package main

import (
	"errors"

	"github.com/tyler-smith/go-bip39"
)

var (
	ErrMnemonicIsInvalid = errors.New("mnemonic is invalid")
)

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return "", ErrMnemonicIsInvalid
	}

	return mnemonic, nil
}
