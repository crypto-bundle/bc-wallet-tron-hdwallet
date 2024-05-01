package wallet_manager

import "errors"

var (
	ErrPassedWalletNotFound                  = errors.New("passed wallet not found")
	ErrUnableCastPluginEntryToPoolUnitWorker = errors.New("unable to cast plugin to pool unit worker")
)
