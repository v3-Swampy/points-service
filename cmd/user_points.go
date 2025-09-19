package cmd

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/v3-Swampy/points-service/cmd/util"
)

type userPointsParams struct {
	Address         string // user address
	TradePoints     uint32 // trade points
	LiquidityPoints uint32 // liquidity points
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

	getUserPointsCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user points",
		Run:   getUserPoints,
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
)

func init() {
	rootCmd.AddCommand(userPointsCmd)

	userPointsCmd.AddCommand(getUserPointsCmd)
	hookUserPointsParams(getUserPointsCmd, false, false)

	userPointsCmd.AddCommand(increaseUserPointsCmd)
	hookUserPointsParams(increaseUserPointsCmd, true, true)

	userPointsCmd.AddCommand(decreaseUserPointsCmd)
	hookUserPointsParams(decreaseUserPointsCmd, true, true)
}

func getUserPoints(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()

	user, err := storeCtx.UserService.Get(pointsParams.Address)
	if err != nil {
		logrus.WithError(err).Info("Failed to get user points")
		return
	}

	logrus.Info("Succeed to get user points:")
	logrus.WithFields(logrus.Fields{
		"address":         user.Address,
		"tradePoints":     user.TradePoints,
		"liquidityPoints": user.LiquidityPoints,
	}).Info("")
}

func incrUserPoints(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()
	deltaUpsertUserPoints(storeCtx, false)
}

func decrUserPoints(cmd *cobra.Command, args []string) {
	storeCtx := util.MustInitStoreContext()
	defer storeCtx.Close()
	deltaUpsertUserPoints(storeCtx, true)
}

func deltaUpsertUserPoints(storeCtx util.StoreContext, decrease bool) {
	if pointsParams.TradePoints == 0 && pointsParams.LiquidityPoints == 0 {
		logrus.Info("At least one of --trade or --liquidity is required.")
		return
	}

	tradePoints := decimal.NewFromInt(int64(pointsParams.TradePoints))
	liquidityPoints := decimal.NewFromInt(int64(pointsParams.LiquidityPoints))
	if decrease {
		tradePoints = tradePoints.Neg()
		liquidityPoints = liquidityPoints.Neg()
	}

	if err := storeCtx.UserService.
		DeltaUpsert(pointsParams.Address, tradePoints, liquidityPoints); err != nil {
		logrus.WithError(err).Info("Failed to update user points")
		return
	}
	logrus.Info("Succeed to update user points")
}

func hookUserPointsParams(cmd *cobra.Command, hookTrade, hookLiquidity bool) {
	cmd.Flags().StringVarP(
		&pointsParams.Address, "user", "u", "", "user address",
	)
	cmd.MarkFlagRequired("user")

	if hookTrade {
		cmd.Flags().Uint32VarP(
			&pointsParams.TradePoints, "trade", "t", 0, "trade points",
		)
	}

	if hookLiquidity {
		cmd.Flags().Uint32VarP(
			&pointsParams.LiquidityPoints, "liquidity", "l", 0, "liquidity points",
		)
	}
}
