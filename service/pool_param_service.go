package service

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/shopspring/decimal"
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

func (service *PoolParamService) GetOrDefault(pool string, defaultValue model.PoolParams) (*model.PoolParams, error) {
	var param model.PoolParams
	found, err := service.store.Get(&param, "address = ?", pool)
	if err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to get pool param values by address")
	}

	if found {
		return &param, nil
	}

	bean := &model.PoolParams{
		Address:         pool,
		TradeWeight:     defaultValue.TradeWeight,
		LiquidityWeight: defaultValue.LiquidityWeight,
	}
	if err := service.store.DB.Create(bean).Error; err != nil {
		return nil, err
	}

	return bean, nil
}

func (service *PoolParamService) Upsert(pool string, tradeWeight, liquidityWeight decimal.Decimal) error {
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
	if tradeWeight.IsPositive() {
		newParam["trade_weight"] = tradeWeight
	}
	if liquidityWeight.IsPositive() {
		newParam["liquidity_weight"] = liquidityWeight
	}

	return service.store.DB.Model(&model.PoolParams{}).
		Where("id = ?", param.ID).
		Updates(newParam).Error
}

func (service *PoolParamService) List() (params []*model.PoolParams, err error) {
	db := service.store.DB.Model(&model.PoolParams{})

	if err = db.Order("id ASC").Find(&params).Error; err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to list pool params")
	}

	return
}

func (service *PoolParamService) MustListPoolAddresses() []string {
	list, err := service.List()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get pools")
	}

	pools := make([]string, 0, len(list))
	for _, pool := range list {
		pools = append(pools, pool.Address)
	}

	return pools
}
