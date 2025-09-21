package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3-Swampy/points-service/cmd/util"
)

type poolWeightParams struct {
	Address         string // pool address
	TradeWeight     uint8  // trade weight
	LiquidityWeight uint8  // liquidity weight
}

var (
	weightParams poolWeightParams

	poolWeightCmd = &cobra.Command{
		Use:   "poolweight",
		Short: "Pool weight utility toolset",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	addPoolWeightCmd = &cobra.Command{
		Use:   "add",
		Short: "Upsert pool weight values",
		Run:   addPoolWeight,
	}

	updatePoolWeightCmd = &cobra.Command{
		Use:   "update",
		Short: "Update pool weight values",
		Run:   updatePoolWeight,
	}

	listPoolWeightCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available pool weight values",
		Run:   listPoolWeight,
	}
)

func init() {
	rootCmd.AddCommand(poolWeightCmd)

	poolWeightCmd.AddCommand(addPoolWeightCmd)
	hookPoolWeightParams(addPoolWeightCmd, true, true)

	poolWeightCmd.AddCommand(updatePoolWeightCmd)
	hookPoolWeightParams(updatePoolWeightCmd, false, false)

	poolWeightCmd.AddCommand(listPoolWeightCmd)
}

func addPoolWeight(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	err := validatePoolWeightParams()
	if err != nil {
		logrus.WithField("config", weightParams).WithError(err).Info("Invalid command config")
		return
	}

	if err := storeCtx.PoolParamService.
		Upsert(weightParams.Address, weightParams.TradeWeight, weightParams.LiquidityWeight); err != nil {
		logrus.WithError(err).Info("Failed to add pool weight values")
		return
	}
	logrus.Info("Succeed to add pool weight values")
}

func updatePoolWeight(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	err := validatePoolWeightParams()
	if err != nil {
		logrus.WithField("config", weightParams).WithError(err).Info("Invalid command config")
		return
	}

	if weightParams.TradeWeight == 0 && weightParams.LiquidityWeight == 0 {
		logrus.Info("At least one of --trade or --liquidity is required.")
		return
	}

	if err := storeCtx.PoolParamService.
		Upsert(weightParams.Address, weightParams.TradeWeight, weightParams.LiquidityWeight); err != nil {
		logrus.WithError(err).Info("Failed to update pool weight values")
		return
	}
	logrus.Info("Succeed to update pool weight values")
}

func listPoolWeight(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	list, err := storeCtx.PoolParamService.List()
	if err != nil {
		logrus.WithError(err).Info("Failed to get pool weight values")
		return
	}

	if len(list) == 0 {
		logrus.Info("No pool weight values found")
		return
	}

	logrus.WithField("total", len(list)).Info("Pool weight values loaded:")
	for i, params := range list {
		logrus.WithFields(logrus.Fields{
			"name":            fmt.Sprintf("%s/%s", params.Token0Symbol, params.Token1Symbol),
			"address":         params.Address,
			"tradeWeight":     params.TradeWeight,
			"liquidityWeight": params.LiquidityWeight,
		}).Info("Pool #", i)
	}
}

func validatePoolWeightParams() error {
	if !common.IsHexAddress(weightParams.Address) {
		return errors.Errorf("Invalid hex address of pool %v", weightParams.Address)
	}
	return nil
}

func hookPoolWeightParams(cmd *cobra.Command, tradeWeightMust, liquidityWeightMust bool) {
	cmd.Flags().StringVarP(
		&weightParams.Address, "pool", "p", "", "pool address",
	)
	cmd.MarkFlagRequired("pool")

	cmd.Flags().Uint8VarP(
		&weightParams.TradeWeight, "trade", "t", 0, "trade weight",
	)
	if tradeWeightMust {
		cmd.MarkFlagRequired("trade")
	}

	cmd.Flags().Uint8VarP(
		&weightParams.LiquidityWeight, "liquidity", "l", 0, "liquidity weight",
	)
	if liquidityWeightMust {
		cmd.MarkFlagRequired("liquidity")
	}
}
