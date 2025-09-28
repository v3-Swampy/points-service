package cmd

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3-Swampy/points-service/cmd/util"
)

type poolWeightParams struct {
	Address              string          // pool address
	TradeWeight          decimal.Decimal // trade weight
	LiquidityWeight      decimal.Decimal // liquidity weight
	TradeWeightParam     string
	LiquidityWeightParam string
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

	getPoolWeightCmd = &cobra.Command{
		Use:   "get",
		Short: "Get pool weight values",
		Run:   getPoolWeight,
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
	hookPoolWeightParams(updatePoolWeightCmd, true, true)

	poolWeightCmd.AddCommand(getPoolWeightCmd)
	hookPoolWeightParams(getPoolWeightCmd, false, false)

	poolWeightCmd.AddCommand(listPoolWeightCmd)
}

func addPoolWeight(cmd *cobra.Command, args []string) {
	upsertPoolWeight()
}

func updatePoolWeight(cmd *cobra.Command, args []string) {
	upsertPoolWeight()
}

func upsertPoolWeight() {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	if err := validatePoolWeightParams(true, true); err != nil {
		logrus.WithError(err).Info("Invalid command config")
		return
	}

	if weightParams.TradeWeight.IsZero() && weightParams.LiquidityWeight.IsZero() {
		logrus.Info("At least one of --trade or --liquidity is required.")
		return
	}

	if err := storeCtx.PoolParamService.
		Upsert(weightParams.Address, weightParams.TradeWeight, weightParams.LiquidityWeight); err != nil {
		logrus.WithError(err).Info("Failed to upsert pool weight values")
		return
	}

	logrus.Info("Succeed to upsert pool weight values")
}

func getPoolWeight(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	if err := validatePoolWeightParams(false, false); err != nil {
		logrus.WithError(err).Info("Invalid command config")
		return
	}

	pool, err := storeCtx.PoolParamService.Get(weightParams.Address)
	if err != nil {
		logrus.WithError(err).Info("Failed to get pool weight values")
		return
	}

	logrus.WithFields(logrus.Fields{
		"address":         pool.Address,
		"tradeWeight":     pool.TradeWeight,
		"liquidityWeight": pool.LiquidityWeight,
	}).Info("Succeed to get pool weight values")
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
			"address":         params.Address,
			"tradeWeight":     params.TradeWeight,
			"liquidityWeight": params.LiquidityWeight,
		}).Info("Pool #", i)
	}
}

func validatePoolWeightParams(validateTradeWeight bool, validateLiquidityWeight bool) error {
	if !common.IsHexAddress(weightParams.Address) {
		return errors.Errorf("Invalid hex address of pool %v", weightParams.Address)
	}

	if validateTradeWeight {
		matched, err := regexp.MatchString(`^(0|[1-9]\d*)(\.\d{1,3})?$`, weightParams.TradeWeightParam)
		if err != nil {
			return errors.Errorf("Invalid trade weight %v", weightParams.TradeWeightParam)
		}
		if !matched {
			return errors.Errorf("Invalid trade weight %v. Only numbers are supported, with a maximum of three decimal", weightParams.TradeWeightParam)
		}

		tradeWeightParam, err := decimal.NewFromString(weightParams.TradeWeightParam)
		if err != nil {
			return err
		}
		weightParams.TradeWeight = tradeWeightParam
	}

	if validateLiquidityWeight {
		matched, err := regexp.MatchString(`^(0|[1-9]\d*)(\.\d{1,3})?$`, weightParams.LiquidityWeightParam)
		if err != nil {
			return errors.Errorf("Invalid liquidity weight %v", weightParams.LiquidityWeightParam)
		}
		if !matched {
			return errors.Errorf("Invalid liquidity weight %v. Only numbers are supported, with a maximum of three decimal", weightParams.LiquidityWeightParam)
		}

		liquidityWeight, err := decimal.NewFromString(weightParams.LiquidityWeightParam)
		if err != nil {
			return err
		}
		weightParams.LiquidityWeight = liquidityWeight
	}

	return nil
}

func hookPoolWeightParams(cmd *cobra.Command, hookTradeWeight, hookLiquidityWeight bool) {
	cmd.Flags().StringVarP(
		&weightParams.Address, "pool", "p", "", "pool address",
	)
	cmd.MarkFlagRequired("pool")

	if hookTradeWeight {
		cmd.Flags().StringVarP(
			&weightParams.TradeWeightParam, "trade", "t", "0", "trade weight",
		)
	}

	if hookLiquidityWeight {
		cmd.Flags().StringVarP(
			&weightParams.LiquidityWeightParam, "liquidity", "l", "0", "liquidity weight",
		)
	}
}
