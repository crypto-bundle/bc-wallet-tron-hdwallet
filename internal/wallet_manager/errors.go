package wallet_manager

import "errors"

var (
	ErrPassedWalletNotFound                  = errors.New("passed wallet not found")
	ErrMethodUnimplemented                   = errors.New("called service method unimplemented")
	ErrUnableCastPluginEntryToPoolUnitMaker  = errors.New("unable to cast plugin to pool unit maker")
	ErrUnableCastPluginEntryToPoolUnitWorker = errors.New("unable to cast plugin to pool unit worker")
)
