package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type unitWrapper struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	Timer      *time.Timer
	Unit       WalletPoolUnitService
}

type Pool struct {
	logger *zap.Logger
	cfg    configService

	runTimeCtx context.Context

	encryptSvc encryptService

	walletUnitsCount uint
	walletUnits      map[uuid.UUID]*unitWrapper
}

func (p *Pool) AddAndStartWalletUnit(_ context.Context,
	walletUUID uuid.UUID,
	timeToLive time.Duration,
	mnemonicEncryptedData []byte,
) error {
	wuWrapper, isExists := p.walletUnits[walletUUID]
	if isExists {
		wuWrapper.Timer.Reset(timeToLive)

		return nil
	}

	unitCtx, cancelFunc := context.WithCancel(p.runTimeCtx)

	started := make(chan struct{})
	timer := time.AfterFunc(timeToLive, func() {
		<-started
		return
	})

	walletUnit := newMnemonicWalletPoolUnit(p.logger, walletUUID, p.encryptSvc, mnemonicEncryptedData)

	wrapper := &unitWrapper{
		Ctx:        unitCtx,
		CancelFunc: cancelFunc,
		Timer:      timer,
		Unit:       walletUnit,
	}
	p.walletUnits[walletUUID] = wrapper
	p.walletUnitsCount++

	err := walletUnit.Run(unitCtx)
	if err != nil {
		return err
	}

	go func(wrapped *unitWrapper) {
		for {
			select {
			case fired, _ := <-wrapped.Timer.C:
				loopErr := wrapped.Unit.Shutdown(wrapped.Ctx)
				if loopErr != nil {
					p.logger.Error("unable to unload wallet data by ticker", zap.Error(loopErr),
						zap.Time(app.TickerEventTriggerTimeTag, fired))
					continue
				}

			case <-wrapped.Ctx.Done():
				loopErr := wrapped.Unit.Shutdown(wrapped.Ctx)
				if loopErr != nil {
					p.logger.Error("unable to shutdown by ctx cancel", zap.Error(loopErr))
				}
			}
		}
	}(wrapper)

	started <- struct{}{}

	return nil
}

func (p *Pool) UnloadWalletUnit(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
) error {
	wUint, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return ErrPassedWalletNotFound
	}

	wUint.CancelFunc()
	wUint.Unit = nil
	p.walletUnits[mnemonicWalletUUID] = nil

	return nil
}

func (p *Pool) GetAddressByPath(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return "", ErrPassedWalletNotFound
	}

	return wUnit.Unit.GetAddressByPath(ctx, account, change, index)
}

func (p *Pool) GetAddressesByPathByRange(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return ErrPassedWalletNotFound
	}

	return wUnit.Unit.GetAddressesByPathByRange(ctx,
		rangeIterable, marshallerCallback)
}

func (p *Pool) LoadAddressByPath(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (*string, error) {
	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, nil
	}

	return wUnit.Unit.LoadAddressByPath(ctx, account, change, index)
}

func (p *Pool) SignData(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	dataForSign []byte,
) (*string, []byte, error) {
	wUnit, isExists := p.walletUnits[mnemonicUUID]
	if !isExists {
		p.logger.Error("wallet is not exists in wallet pool",
			zap.String(app.WalletUUIDTag, mnemonicUUID.String()))

		return nil, nil, ErrPassedWalletNotFound
	}

	return wUnit.Unit.SignData(ctx, account, change, index, dataForSign)
}

func newWalletPool(logger *zap.Logger,
	cfg configService,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		logger:           logger,
		cfg:              cfg,
		encryptSvc:       encryptSrv,
		walletUnits:      make(map[uuid.UUID]*unitWrapper),
		walletUnitsCount: 0,
	}
}
