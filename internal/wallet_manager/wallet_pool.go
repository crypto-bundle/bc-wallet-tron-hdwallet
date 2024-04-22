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
	TTL        time.Duration
	Unit       WalletPoolUnitService
	OnShutDown func(walletUUID uuid.UUID)
}

func (w *unitWrapper) Shutdown() error {
	err := w.Unit.Shutdown(w.Ctx)
	if err != nil {
		return err
	}

	w.Unit = nil
	w.Timer.Stop()
	w.Timer = nil
	w.Ctx = nil

	return nil
}

type Pool struct {
	logger *zap.Logger
	cfg    configService

	runTimeCtx context.Context

	encryptSvc encryptService

	walletUnits map[uuid.UUID]*unitWrapper
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

	walletUnit := newMnemonicWalletPoolUnit(p.logger, walletUUID, p.encryptSvc, mnemonicEncryptedData)

	wrapper := &unitWrapper{
		Ctx:        unitCtx,
		CancelFunc: cancelFunc,
		Timer:      nil, // will be filled in go-routine
		TTL:        timeToLive,
		Unit:       walletUnit,
		OnShutDown: p.unloadWalletUnit,
	}
	p.walletUnits[walletUUID] = wrapper

	err := walletUnit.Run(unitCtx)
	if err != nil {
		return err
	}
	started := make(chan struct{})

	go func(wrapped *unitWrapper) {
		walletUUIDInt := wrapped.Unit.GetMnemonicUUID()
		wrapped.Timer = time.NewTimer(wrapped.TTL)
		started <- struct{}{}

		for {
			select {
			case fired, _ := <-wrapped.Timer.C:
				loopErr := wrapped.Shutdown()
				if loopErr != nil {
					p.logger.Error("unable to unload wallet data by ticker", zap.Error(loopErr),
						zap.Time(app.TickerEventTriggerTimeTag, fired))
				}

				wrapped.OnShutDown(*walletUUIDInt)

				return

			case <-wrapped.Ctx.Done():
				loopErr := wrapped.Shutdown()
				if loopErr != nil {
					p.logger.Error("unable to shutdown by ctx cancel", zap.Error(loopErr))
				}

				return
			}
		}
	}(wrapper)

	<-started
	close(started)

	return nil
}

func (p *Pool) UnloadWalletUnit(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
) (*uuid.UUID, error) {
	wUint, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, nil
	}
	walletUUID := wUint.Unit.GetMnemonicUUID()

	wUint.CancelFunc()

	p.unloadWalletUnit(mnemonicWalletUUID)

	return walletUUID, nil
}

func (p *Pool) unloadWalletUnit(mnemonicWalletUUID uuid.UUID) {
	p.walletUnits[mnemonicWalletUUID] = nil
	delete(p.walletUnits, mnemonicWalletUUID)

	return
}

func (p *Pool) UnloadMultipleWalletUnit(ctx context.Context,
	mnemonicWalletUUIDs []uuid.UUID,
) error {
	for _, v := range mnemonicWalletUUIDs {
		wUint, isExists := p.walletUnits[v]
		if !isExists {
			continue
		}

		wUint.CancelFunc()

		p.walletUnits[v] = nil
		delete(p.walletUnits, v)
	}

	return nil
}

func (p *Pool) GetAddressByPath(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (*string, error) {
	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, ErrPassedWalletNotFound
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

func NewWalletPool(ctx context.Context,
	logger *zap.Logger,
	cfg configService,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		runTimeCtx:  ctx,
		logger:      logger,
		cfg:         cfg,
		encryptSvc:  encryptSrv,
		walletUnits: make(map[uuid.UUID]*unitWrapper),
	}
}
