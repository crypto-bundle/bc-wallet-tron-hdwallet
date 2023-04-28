package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/app"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/config"
	grpcHandlers "gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/grpc"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/mnemonic"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/wallet_data"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/internal/wallet_manager"
	"gitlab.heronodes.io/bc-platform/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api"

	commonHealthcheck "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-redis/pkg/redis"

	_ "github.com/mailru/easyjson/gen"
	"go.uber.org/zap"
)

// DO NOT EDIT THIS VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// Version - version time.RFC3339.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	Version = "DEVELOPMENT.VESION"

	// ReleaseTag - release tag in TAG.%Y-%m-%dT%H-%M-%SZ.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "DEVELOPMENT.RELEASE_TAG"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "DEVELOPMENT.COMMIT_HASH"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "DEVELOPMENT.SHORT_COMMIT_HASH"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber uint64 = 0

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS uint64 = 0
)

const ApplicationName = "bc-wallet-tron-hdwallet-api"

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, secretSrv, err := config.Prepare(ctx, Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS, ApplicationName)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main")

	pgConn := commonPostgres.NewConnection(context.Background(), appCfg, loggerEntry)
	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	natsConnSvc := commonNats.NewConnection(ctx, appCfg, loggerEntry)
	err = natsConnSvc.Connect()
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	redisSvc := commonRedis.NewConnection(ctx, appCfg, loggerEntry)
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	redisConn, err := redisSvc.Connect(ctx)
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}
	redisClient := redisConn.GetClient()

	walletDataSrv := wallet_data.NewService(loggerEntry, pgConn)
	mnemonicWalletDataSrv, err := mnemonic_wallet_data.NewService(loggerEntry, appCfg,
		pgConn, redisClient, natsConnSvc)
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	mnemonicGenerator := mnemonic.NewMnemonicGenerator(loggerEntry,
		appCfg.GetDefaultMnemonicWordsCount())

	walletService, err := wallet_manager.NewService(loggerEntry, appCfg, secretSrv,
		walletDataSrv, mnemonicWalletDataSrv,
		pgConn, mnemonicGenerator)
	if err != nil {
		loggerEntry.Fatal("unable to create wallet service instance", zap.Error(err))
	}

	apiHandlers, err := grpcHandlers.New(ctx, loggerEntry, walletService)
	if err != nil {
		loggerEntry.Fatal("unable to init grpc handlers", zap.Error(err))
	}

	srv, err := hdwallet_api.NewServer(ctx, loggerEntry, appCfg, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = walletService.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init wallet service", zap.Error(err))
	}

	err = walletService.Run(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to run wallet service", zap.Error(err))
	}

	err = srv.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	commonHealthcheck.NewHTTPHealthChecker(loggerEntry)

	go func() {
		err = srv.ListenAndServe(ctx)
		if err != nil {
			loggerEntry.Error("unable to start grpc", zap.Error(err),
				zap.String("port", appCfg.GetBindPort()))
		}
	}()

	loggerEntry.Info("application started successfully", zap.String(app.GRPCBindPortTag, appCfg.GetBindPort()))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	loggerEntry.Warn("shutdown application")
	cancelCtxFunc()

	walletShutdownErr := walletService.Shutdown(ctx)
	if walletShutdownErr != nil {
		log.Printf("%s:%e", walletShutdownErr.Error(), walletShutdownErr)
	}

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
