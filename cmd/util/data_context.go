package util

import (
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/service"
)

type StoreContext struct {
	Store            *store.Store
	PoolParamService *service.PoolParamService
	UserService      *service.UserService
}

func MustInitStoreContext() StoreContext {
	var ctx StoreContext

	// init database
	storeConfig := store.MustNewConfigFromViper()
	db := storeConfig.MustOpenOrCreate(model.Tables...)
	ctx.Store = store.NewStore(db)

	// init services
	ctx.PoolParamService = service.NewPoolParamService(ctx.Store)
	ctx.UserService = service.NewUserService(ctx.Store)

	return ctx
}

func (ctx *StoreContext) Close() {
	if ctx.Store != nil {
		ctx.Store.Close()
	}
}
