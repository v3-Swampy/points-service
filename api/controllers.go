package api

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/gin-gonic/gin"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/service"
)

type Controller struct {
	userService *service.UserService
	poolService *service.PoolService
}

func NewController(store *store.Store) *Controller {
	return &Controller{
		userService: service.NewUserService(store),
		poolService: service.NewPoolService(store),
	}
}

func (controller *Controller) listUsers(c *gin.Context) (any, error) {
	var input model.UserPagingRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, api.ErrValidation(err)
	}

	return controller.userService.List(input)
}

func (controller *Controller) listPools(c *gin.Context) (any, error) {
	var input model.PoolPagingRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, api.ErrValidation(err)
	}

	return controller.poolService.List(input)
}
