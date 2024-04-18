package grpc

import "errors"

var (
	ErrMissedRequiredData     = errors.New("missed required data")
	ErrWalletUUIDMismatched   = errors.New("wallet uuid mismatched")
	ErrDataIsNotValid         = errors.New("received data is not valid")
	ErrUnableReadGrpcMetadata = errors.New("unable to read grpc metadata")
)
