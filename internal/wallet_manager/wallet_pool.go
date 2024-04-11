package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.Logger
	cfg    configService

	runTimeCtx context.Context

	encryptSrv encryptService

	walletUnitsCount uint

	walletUnitTimers map[uuid.UUID]*time.Timer

	walletUnits map[uuid.UUID]WalletPoolUnitService
}

func (p *Pool) Init(ctx context.Context) error {
	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Init(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) Run(ctx context.Context) error {
	p.runTimeCtx = ctx

	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Run(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) Shutdown(ctx context.Context) error {
	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Shutdown(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) SetWalletUnits(ctx context.Context,
	walletUnits map[uuid.UUID]WalletPoolUnitService,
) error {
	if len(p.walletUnits) > 0 {
		return ErrWalletPoolIsNotEmpty
	}

	if len(walletUnits) == 0 {
		return ErrPassedWalletPoolUnitIsEmpty
	}

	p.walletUnits = walletUnits
	p.walletUnitsCount = uint(len(walletUnits))

	return nil
}

func (p *Pool) AddAWalletUnit(ctx context.Context,
	walletUUID uuid.UUID,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	p.walletUnits[walletUUID] = walletUnit
	p.walletUnitsCount++

	return nil
}

func (p *Pool) AddAndStartWalletUnit(_ context.Context,
	walletUUID uuid.UUID,
	timeToLive time.Duration,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	err := walletUnit.Init(p.runTimeCtx)
	if err != nil {
		return err
	}

	err = walletUnit.Run(p.runTimeCtx)
	if err != nil {
		return err
	}

	p.walletUnits[walletUUID] = walletUnit
	p.walletUnitsCount++

	timer, isExists := p.walletUnitTimers[walletUUID]
	if isExists {
		timer.Reset(timeToLive)

		return nil
	}

	timer = time.NewTimer(timeToLive)
	p.walletUnitTimers[walletUUID] = timer

	go func() {
		for {
			select {
			case fired, _ := <-timer.C:
				loopErr := walletUnit.Shutdown(p.runTimeCtx)
				if loopErr != nil {
					p.logger.Error("unable to unload wallet data by ticker", zap.Error(err),
						zap.Time(app.TickerEventTriggerTimeTag, fired))
					continue
				}

			case <-p.runTimeCtx.Done():
				loopErr := walletUnit.Shutdown(p.runTimeCtx)
				if loopErr != nil {
					p.logger.Error("unable to shutdown by ctx cancel", zap.Error(err))
				}
			}
		}
	}()

	return nil
}

func (p *Pool) GetAddressByPath(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		return "", ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressByPath(ctx, mnemonicWalletUUID, account, change, index)
}

func (p *Pool) GetAddressesByPathByRange(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		return ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressesByPathByRange(ctx, mnemonicWalletUUID,
		rangeIterable, marshallerCallback)
}

func (p *Pool) SignData(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transactionData []byte,
) (*string, []byte, error) {
	poolUnit, isExists := p.walletUnits[mnemonicUUID]
	if !isExists {
		p.logger.Error("wallet is not exists in wallet pool",
			zap.String(app.WalletUUIDTag, mnemonicUUID.String()))

		return nil, nil, ErrPassedWalletNotFound
	}

	return poolUnit.Sign(ctx, account, change, index, transactionData)
}

func newWalletPool(logger *zap.Logger,
	cfg configService,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		logger:           logger,
		cfg:              cfg,
		encryptSrv:       encryptSrv,
		walletUnits:      make(map[uuid.UUID]WalletPoolUnitService, 0),
		walletUnitsCount: 0,
	}
}
