package grpc

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameValidateMnemonic = "ValidateMnemonic"
)

type validateMnemonicHandler struct {
	l *zap.Logger

	appEncryptorSvc      encryptService
	mnemonicValidatorSvc mnemonicValidatorService
}

// nolint:funlen // fixme
func (h *validateMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.ValidateMnemonicRequest,
) (*pbApi.ValidateMnemonicResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &ValidateMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	decryptedData, err := h.appEncryptorSvc.Decrypt(vf.EncryptedMnemonicData)
	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	defer func() {
		for i := range decryptedData {
			decryptedData[i] = 0
		}
	}()

	if decryptedData == nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	isValid := h.mnemonicValidatorSvc.IsMnemonicValid(string(decryptedData))
	if !isValid {
		return nil, status.Error(codes.InvalidArgument, "seed phrase is not valid")
	}

	mnemonicHash := fmt.Sprintf("%x", sha256.Sum256(decryptedData))
	req.MnemonicIdentity.WalletHash = mnemonicHash

	return &pbApi.ValidateMnemonicResponse{
		MnemonicIdentity: req.MnemonicIdentity,
		IsValid:          isValid,
	}, nil
}

func MakeValidateMnemonicHandler(loggerEntry *zap.Logger,
	appEncryptorSvc encryptService,
	mnemoValidatorSvc mnemonicValidatorService,
) *validateMnemonicHandler {
	return &validateMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameValidateMnemonic)),

		appEncryptorSvc:      appEncryptorSvc,
		mnemonicValidatorSvc: mnemoValidatorSvc,
	}
}
