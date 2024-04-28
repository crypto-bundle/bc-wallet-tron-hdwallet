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
	MethodNameGenerateMnemonic = "GenerateMnemonic"
)

type generateMnemonicHandler struct {
	l *zap.Logger

	generatorSvc    mnemonicGeneratorService
	appEncryptorSvc encryptService
}

// nolint:funlen // fixme
func (h *generateMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	var err error
	tCtx, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	vf := &GenerateMnemonicForm{}
	valid, err := vf.LoadAndValidate(tCtx, req)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	seedPhrase, err := h.generatorSvc.Generate(ctx)
	if err != nil {
		h.l.Error("unable to generate mnemonic", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, vf.WalletUUID))

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	seedPhraseRaw := []byte(seedPhrase)

	encryptedMnemonicData, err := h.appEncryptorSvc.Encrypt(seedPhraseRaw)
	if err != nil {
		return nil, err
	}

	mnemonicHash := fmt.Sprintf("%x", sha256.Sum256(seedPhraseRaw))

	req.MnemonicIdentity.WalletHash = mnemonicHash

	defer func() {
		for i := range seedPhraseRaw {
			seedPhraseRaw[i] = 0
		}
		seedPhraseRaw = nil
		seedPhrase = ""
	}()

	return &pbApi.GenerateMnemonicResponse{
		MnemonicIdentity:      req.MnemonicIdentity,
		EncryptedMnemonicData: encryptedMnemonicData,
	}, nil
}

func MakeGenerateMnemonicHandler(loggerEntry *zap.Logger,
	generatorSvc mnemonicGeneratorService,
	appEncryptorSvc encryptService,
) *generateMnemonicHandler {
	return &generateMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGenerateMnemonic)),

		generatorSvc:    generatorSvc,
		appEncryptorSvc: appEncryptorSvc,
	}
}
