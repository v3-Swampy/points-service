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

type userPointsParams struct {
	Address              string          // user address
	TradePoints          decimal.Decimal // trade points
	LiquidityPoints      decimal.Decimal // liquidity points
	TradePointsParam     uint32
	LiquidityPointsParam string
}

var (
	pointsParams userPointsParams

	userPointsCmd = &cobra.Command{
		Use:   "userpoints",
		Short: "User points utility toolset",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	insertUserPointsCmd = &cobra.Command{
		Use:   "insert",
		Short: "Insert user points",
		Run:   insertUserPoints,
	}

	increaseUserPointsCmd = &cobra.Command{
		Use:   "increase",
		Short: "Increase user points",
		Run:   incrUserPoints,
	}

	decreaseUserPointsCmd = &cobra.Command{
		Use:   "decrease",
		Short: "Decrease user points",
		Run:   decrUserPoints,
	}

	getUserPointsCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user points",
		Run:   getUserPoints,
	}
)

func init() {
	rootCmd.AddCommand(userPointsCmd)

	userPointsCmd.AddCommand(insertUserPointsCmd)
	hookUserPointsParams(insertUserPointsCmd, true, true)

	userPointsCmd.AddCommand(increaseUserPointsCmd)
	hookUserPointsParams(increaseUserPointsCmd, true, true)

	userPointsCmd.AddCommand(decreaseUserPointsCmd)
	hookUserPointsParams(decreaseUserPointsCmd, true, true)

	userPointsCmd.AddCommand(getUserPointsCmd)
	hookUserPointsParams(getUserPointsCmd, false, false)
}

func insertUserPoints(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	err := validateAndConvertUserPointsParams(true, true)
	if err != nil {
		logrus.WithError(err).Info("Invalid command config")
		return
	}

	if _, err = storeCtx.UserService.
		Add(pointsParams.Address, pointsParams.TradePoints, pointsParams.LiquidityPoints); err != nil {
		logrus.WithError(err).Info("Failed to insert user points")
		return
	}

	logrus.Info("Succeed to insert user points")
}

func incrUserPoints(cmd *cobra.Command, args []string) {
	deltaUpdateUserPoints(false)
}

func decrUserPoints(cmd *cobra.Command, args []string) {
	deltaUpdateUserPoints(true)
}

func deltaUpdateUserPoints(decrease bool) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	err := validateAndConvertUserPointsParams(true, true)
	if err != nil {
		logrus.WithError(err).Info("Invalid command config")
		return
	}

	if pointsParams.TradePoints.Equal(decimal.Zero) && pointsParams.LiquidityPoints.Equal(decimal.Zero) {
		logrus.Info("At least one of --trade or --liquidity is required.")
		return
	}

	if decrease {
		pointsParams.TradePoints = pointsParams.TradePoints.Neg()
		pointsParams.LiquidityPoints = pointsParams.LiquidityPoints.Neg()
	}

	if err := storeCtx.UserService.
		DeltaUpdate(pointsParams.Address, pointsParams.TradePoints, pointsParams.LiquidityPoints); err != nil {
		logrus.WithError(err).Info("Failed to update user points")
		return
	}

	logrus.Info("Succeed to update user points")
}

func getUserPoints(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	err := validateAndConvertUserPointsParams(false, false)
	if err != nil {
		logrus.WithError(err).Info("Invalid command config")
		return
	}

	user, err := storeCtx.UserService.Get(pointsParams.Address)
	if err != nil {
		logrus.WithError(err).Info("Failed to get user points")
		return
	}

	logrus.WithFields(logrus.Fields{
		"address":         user.Address,
		"tradePoints":     user.TradePoints,
		"liquidityPoints": user.LiquidityPoints,
	}).Info("Succeed to get user points")
}

func validateAndConvertUserPointsParams(validateTradePoints bool, validateLiquidityPoints bool) error {
	if !common.IsHexAddress(pointsParams.Address) {
		return errors.Errorf("Invalid hex address of user %v", pointsParams.Address)
	}

	if validateTradePoints {
		pointsParams.TradePoints = decimal.NewFromInt(int64(pointsParams.TradePointsParam))
	}

	if validateLiquidityPoints {
		matched, err := regexp.MatchString(`^(0|[1-9]\d*)(\.\d)?$`, pointsParams.LiquidityPointsParam)
		if err != nil {
			return errors.Errorf("Invalid liquidity points value %v", pointsParams.LiquidityPointsParam)
		}
		if !matched {
			return errors.Errorf("Invalid liquidity points value %v. Only numbers are supported, with a maximum of one decimal", pointsParams.LiquidityPointsParam)
		}

		liquidityPoints, err := decimal.NewFromString(pointsParams.LiquidityPointsParam)
		if err != nil {
			return err
		}
		pointsParams.LiquidityPoints = liquidityPoints
	}

	return nil
}

func hookUserPointsParams(cmd *cobra.Command, hookTrade, hookLiquidity bool) {
	cmd.Flags().StringVarP(
		&pointsParams.Address, "user", "u", "", "user address",
	)
	cmd.MarkFlagRequired("user")

	if hookTrade {
		cmd.Flags().Uint32VarP(
			&pointsParams.TradePointsParam, "trade", "t", 0, "trade points",
		)
	}

	if hookLiquidity {
		cmd.Flags().StringVarP(
			&pointsParams.LiquidityPointsParam, "liquidity", "l", "0", "liquidity points",
		)
	}
}
