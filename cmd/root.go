package cmd

import (
	"context"
	"sync"

	"github.com/Conflux-Chain/go-conflux-util/cmd"
	"github.com/Conflux-Chain/go-conflux-util/config"
	"github.com/Conflux-Chain/go-conflux-util/log"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/Conflux-Chain/go-conflux-util/viper"
	"github.com/openweb3/web3go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3-Swampy/points-service/api"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/service"
	"github.com/v3-Swampy/points-service/sync/parsing"
)

var rootCmd = &cobra.Command{
	Use: "points-service",
	Run: start,
}

func init() {
	cobra.OnInitialize(func() {
		config.MustInit("PS")
	})

	log.BindFlags(rootCmd)
}

func start(*cobra.Command, []string) {
	logrus.Info("Starting service ...")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// init database
	storeConfig := store.MustNewConfigFromViper()
	db := storeConfig.MustOpenOrCreate(model.Tables...)
	store := store.NewStore(db)

	// init blockchain client
	var blockchainConfig struct {
		URL string
	}
	viper.MustUnmarshalKey("blockchain", &blockchainConfig)
	client, err := web3go.NewClient(blockchainConfig.URL)
	cmd.FatalIfErr(err, "Failed to create blockchain client")
	defer client.Close()

	// init stat services
	var swappiConfig service.SwappiConfig
	viper.MustUnmarshalKey("sync.swappi", &swappiConfig)
	statService := service.NewStatService(swappiConfig, client, store)

	// init pools
	paramsService := service.NewPoolParamService(store)
	pools := paramsService.MustListPoolAddresses()

	// start sync service
	var syncConfig parsing.Config
	viper.MustUnmarshalKey("sync", &syncConfig)
	syncConfig.Pools = pools
	syncService, err := parsing.NewService(syncConfig, statService, client)
	cmd.FatalIfErr(err, "Failed to create sync service")
	wg.Add(1)
	go syncService.Run(ctx, &wg)

	// start api
	go api.MustServeFromViper(store)

	logrus.Info("Service started")

	cmd.GracefulShutdown(&wg, cancel)
}

// Execute is the command line entrypoint.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute command")
	}
}
