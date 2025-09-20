package service

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/model"
)

type PoolParamService struct {
	store *store.Store
}

func NewPoolParamService(store *store.Store) *PoolParamService {
	return &PoolParamService{
		store: store,
	}
}

func (service *PoolParamService) Get(pool string) (*model.PoolParams, error) {
	var param model.PoolParams
	found, err := service.store.Get(&param, "address = ?", pool)
	if err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to get pool param values by address")
	}

	if !found {
		return nil, api.ErrValidationStr("Failed to find pool param values by address")
	}

	return &param, nil
}

func (service *PoolParamService) Upsert(pool string, tradeWeight, liquidityWeight uint8) error {
	var param model.PoolParams
	found, err := service.store.Get(&param, "address = ?", pool)
	if err != nil {
		return api.ErrDatabaseCause(err, "Failed to get pool param values by address")
	}

	if !found {
		bean := &model.PoolParams{
			Address:         pool,
			TradeWeight:     tradeWeight,
			LiquidityWeight: liquidityWeight,
		}
		return service.store.DB.Create(bean).Error
	}

	newParam := map[string]any{}
	if tradeWeight > 0 {
		newParam["trade_weight"] = tradeWeight
	}
	if liquidityWeight > 0 {
		newParam["liquidity_weight"] = liquidityWeight
	}

	return service.store.DB.Model(&model.PoolParams{}).
		Where("id = ?", param.ID).
		Updates(newParam).Error
}

func (service *PoolParamService) List() (params []*model.PoolParamInfo, err error) {
	db := service.store.DB.Model(&model.PoolParams{})

	db = db.Select("pool_params.*, pools.token0_symbol, pools.token1_symbol").
		Joins("INNER JOIN pools ON pool_params.address = pools.address").
		Order("id ASC")

	if err = db.Find(&params).Error; err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to list pool params")
	}

	return
}

func (service *PoolParamService) MustListPoolAddresses() []string {
	list, err := service.List()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get pools")
	}

	pools := make([]string, len(list))
	for _, pool := range list {
		pools = append(pools, pool.Address)
	}

	return pools
}
