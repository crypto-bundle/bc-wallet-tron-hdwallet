package grpc

import "errors"

var (
	ErrMissedRequiredData   = errors.New("missed required data")
	ErrWalletUUIDMismatched = errors.New("wallet uuid mismatched")
)
