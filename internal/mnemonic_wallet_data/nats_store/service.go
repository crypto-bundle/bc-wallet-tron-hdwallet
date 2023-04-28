package nats_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"

	tracer "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/app"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/entities"

	originNats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	NatsMnemonicWalletBucketName = "mnemonic-wallets"
)

var (
	ErrPassedNilNatsConn = errors.New("passed nil nats connection")
)

type natsStore struct {
	natsConn originNats.JetStreamContext
	logger   *zap.Logger

	bucketName string
	kvBucket   originNats.KeyValue
}

func (s *natsStore) SetMnemonicWalletItem(ctx context.Context,
	walletItem *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletUUIDTag, walletItem.UUID)

	rawJSON, err := walletItem.MarshalJSON()
	if err != nil {
		return nil, err
	}

	keyUUID := walletItem.UUID.String()
	_, err = s.kvBucket.Put(keyUUID, rawJSON)
	if err != nil {
		return nil, err
	}

	keyHash := walletItem.MnemonicHash
	_, err = s.kvBucket.Put(keyHash, rawJSON)
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func (s *natsStore) GetMnemonicWalletItemByUUID(ctx context.Context,
	MnemonicWalletUUID uuid.UUID,
) (*entities.MnemonicWallet, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletUUIDTag, MnemonicWalletUUID)

	kvEntry, err := s.kvBucket.Get(MnemonicWalletUUID.String())
	if err != nil {
		if errors.Is(err, originNats.ErrKeyNotFound) {
			err = nil // for tracer defer
			return nil, nil
		}
		return nil, err
	}

	rawJSON := kvEntry.Value()
	if rawJSON == nil {
		return nil, nil
	}

	walletItem := &entities.MnemonicWallet{}
	err = walletItem.UnmarshalJSON(rawJSON)
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func (s *natsStore) GetMnemonicWalletItemByHash(ctx context.Context,
	MnemonicWalletHash string,
) (*entities.MnemonicWallet, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.MnemonicWalletHashTag, MnemonicWalletHash)

	kvEntry, err := s.kvBucket.Get(MnemonicWalletHash)
	if err != nil {
		if errors.Is(err, originNats.ErrKeyNotFound) {
			err = nil // for tracer defer
			return nil, nil
		}

		return nil, err
	}

	rawJSON := kvEntry.Value()
	if rawJSON == nil {
		return nil, nil
	}

	walletItem := &entities.MnemonicWallet{}
	err = walletItem.UnmarshalJSON(rawJSON)
	if err != nil {
		return nil, err
	}

	return walletItem, nil
}

func NewNatsStore(logger *zap.Logger,
	cfg configurationService,
	natsConn originNats.JetStreamContext,
) (*natsStore, error) {
	bucketName := strings.ToUpper(fmt.Sprintf("%s__%s__%s", cfg.GetStageName(), app.ApplicationName,
		NatsMnemonicWalletBucketName),
	)

	if natsConn == nil {
		return nil, ErrPassedNilNatsConn
	}

	kvBucket, err := natsConn.KeyValue(bucketName)
	if err != nil {
		if errors.Is(err, originNats.ErrBucketNotFound) {
			logger.Error("nats kv bucket not found. plz create it", zap.Error(err))
		}

		return nil, err
	}

	return &natsStore{
		natsConn:   natsConn,
		logger:     logger,
		kvBucket:   kvBucket,
		bucketName: bucketName,
	}, nil
}
