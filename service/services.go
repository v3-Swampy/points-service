package service

import (
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/blockchain"
)

type Services struct {
	Config    *ConfigService
	PoolParam *PoolParamService
	Pool      *PoolService
	User      *UserService
	Stat      *StatService
}

func NewServices(store *store.Store, vswap *blockchain.Vswap) Services {
	return Services{
		Config:    NewConfigService(store),
		PoolParam: NewPoolParamService(store),
		Pool:      NewPoolService(store),
		User:      NewUserService(store),
		Stat:      NewStatService(store, vswap),
	}
}
