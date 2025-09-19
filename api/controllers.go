package api

import (
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/gin-gonic/gin"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/service"
)

type Controller struct {
	userService   *service.UserService
	poolService   *service.PoolService
	configService *service.ConfigService
}

func NewController(store *store.Store) *Controller {
	return &Controller{
		userService:   service.NewUserService(store),
		poolService:   service.NewPoolService(store),
		configService: service.NewConfigService(store),
	}
}

func (controller *Controller) listUsers(c *gin.Context) (any, error) {
	var input model.UserPagingRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, api.ErrValidation(err)
	}

	total, list, err := controller.userService.List(input)
	if err != nil {
		return nil, err
	}

	users := make([]model.UserInfo, 0)
	for _, u := range list {
		user := model.UserInfo{
			Address:         u.Address,
			TradePoints:     u.TradePoints,
			LiquidityPoints: u.LiquidityPoints,
		}
		users = append(users, user)
	}

	updatedAt, err := controller.configService.GetLastStatPointsTime()
	if err != nil {
		return nil, err
	}

	return model.PagingResultWithUpdatedAt[model.UserInfo]{
		Total:     total,
		Items:     users,
		UpdatedAt: updatedAt,
	}, nil
}

func (controller *Controller) listPools(c *gin.Context) (any, error) {
	var input model.PoolPagingRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, api.ErrValidation(err)
	}

	total, list, err := controller.poolService.List(input)
	if err != nil {
		return nil, err
	}

	pools := make([]model.PoolInfo, 0)
	for _, p := range list {
		pool := model.PoolInfo{
			Address:         p.Address,
			Name:            fmt.Sprintf("%s/%s", p.Token0Symbol, p.Token1Symbol),
			Tvl:             p.Tvl,
			TradePoints:     p.TradePoints,
			LiquidityPoints: p.LiquidityPoints,
		}
		pools = append(pools, pool)
	}

	return model.PagingResult[model.PoolInfo]{
		Total: total,
		Items: pools,
	}, nil
}
