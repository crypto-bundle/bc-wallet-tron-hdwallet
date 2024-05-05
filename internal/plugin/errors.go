package plugin

import "errors"

var (
	ErrUnableCastPluginEntry                 = errors.New("unable to cast plugin entry to traget")
	ErrUnableCastPluginEntryToPoolUnitWorker = errors.New("unable to cast plugin to pool unit worker")
)
