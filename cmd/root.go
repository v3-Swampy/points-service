package cmd

import (
	"context"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/model"
	"sync"

	"github.com/Conflux-Chain/go-conflux-util/cmd"
	"github.com/Conflux-Chain/go-conflux-util/config"
	"github.com/Conflux-Chain/go-conflux-util/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3-Swampy/points-service/api"
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

	_, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	storeConfig := store.MustNewConfigFromViper()
	db := storeConfig.MustOpenOrCreate(model.Tables...)

	go api.MustServeFromViper(db)

	cmd.GracefulShutdown(&wg, cancel)
}

// Execute is the command line entrypoint.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute command")
	}
}
