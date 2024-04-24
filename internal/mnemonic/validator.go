package mnemonic

import (
	"github.com/tyler-smith/go-bip39"
)

type Validator struct {
}

func (v Validator) IsMnemonicValid(mnemonic string) bool {
	return Validate(mnemonic)
}

func NewMnemonicValidator() *Validator {
	return &Validator{}
}

func Validate(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}
