package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.Logger
	cfg    configService

	walletsDataSrv         walletsDataService
	mnemonicWalletsDataSrv mnemonicWalletsDataService
	encryptSrv             encryptService

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

	p.walletUnits = walletUnits

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

	return nil
}

func (p *Pool) AddAndStartWalletUnit(ctx context.Context,
	walletUUID uuid.UUID,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	p.walletUnits[walletUUID] = walletUnit

	err := walletUnit.Init(ctx)
	if err != nil {
		return err
	}

	err = walletUnit.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) GetAddressesByPathByRange(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if isExists {
		return nil, ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressesByPathByRange(ctx, mnemonicWalletUUID,
		accountIndex, internalIndex,
		addressIndexFrom, addressIndexTo)
}

func NewWalletPool(logger *zap.Logger,
	cfg configService,
	walletsDataSrv walletsDataService,
	mnemonicWalletsDataSrv mnemonicWalletsDataService,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		logger:                 logger,
		cfg:                    cfg,
		walletsDataSrv:         walletsDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletsDataSrv,
		encryptSrv:             encryptSrv,
		walletUnits:            make(map[uuid.UUID]WalletPoolUnitService, 0),
	}
}