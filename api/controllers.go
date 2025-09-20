package api

import (
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

// listUsers returns users in pagination view.
//
//	@Summary		List users
//	@Description	List users in pagination view.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			offset		query		int																		false	"The number of skipped records, usually it's pageSize * (pageNumber - 1)"	minimum(0)				default(0)
//	@Param			limit		query		int																		true	"The number of records displayed on the page"								minimum(1)				maximum(100)
//	@Param			sort		query		string																	false	"Sort in ASC or DESC order by sortField"									Enums(asc, desc)		default(desc)
//	@Param			sortField	query		string																	false	"The field used for sorting. The value is trade or liquidity"				Enums(trade, liquidity)	default(trade)
//	@Success		200			{object}	api.BusinessError{data=model.PagingResultWithUpdatedAt[model.UserInfo]}	"Paged users"
//	@Failure		600			{object}	api.BusinessError{data=string}											"Internal server error"
//	@Router			/users		[get]
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

// listPools returns pools in pagination view.
//
//	@Summary		List pools
//	@Description	List pools in pagination view.
//	@Tags			Pool
//	@Accept			json
//	@Produce		json
//	@Param			offset		query		int															false	"The number of skipped records, usually it's pageSize * (pageNumber - 1)"	minimum(0)						default(0)
//	@Param			limit		query		int															true	"The number of records displayed on the page"								minimum(1)						maximum(100)
//	@Param			sort		query		string														false	"Sort in ASC or DESC order by sortField"									Enums(asc, desc)				default(desc)
//	@Param			sortField	query		string														false	"The field used for sorting. The value is tvl, trade or liquidity"			Enums(tvl, trade, liquidity)	default(tvl)
//	@Success		200			{object}	api.BusinessError{data=model.PagingResult[model.PoolInfo]}	"Paged pools"
//	@Failure		600			{object}	api.BusinessError{data=string}								"Internal server error"
//	@Router			/pools		[get]
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
			PoolParamInfo: model.PoolParamInfo{
				Address:         p.Address,
				Token0Symbol:    p.Token0Symbol,
				Token1Symbol:    p.Token1Symbol,
				TradeWeight:     p.TradeWeight,
				LiquidityWeight: p.LiquidityWeight,
			},
			Tvl: p.Tvl,
		}
		pools = append(pools, pool)
	}

	return model.PagingResult[model.PoolInfo]{
		Total: total,
		Items: pools,
	}, nil
}
