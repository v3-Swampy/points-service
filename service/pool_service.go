package service

import (
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/model"
	"gorm.io/gorm"
)

type PoolService struct {
	store *store.Store
}

func NewPoolService(store *store.Store) *PoolService {
	return &PoolService{
		store: store,
	}
}

func (service *PoolService) BatchDeltaUpsert(pools []*model.Pool, dbTx ...*gorm.DB) error {
	db := service.store.DB
	if len(dbTx) > 0 {
		db = dbTx[0]
	}

	var placeholders string
	var params []interface{}
	size := len(pools)
	for i, p := range pools {
		placeholders += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		if i != size-1 {
			placeholders += ",\n\t\t\t"
		}
		params = append(params, []interface{}{
			p.Address, p.Token0, p.Token1, p.Tvl, p.TradePoints, p.LiquidityPoints,
			p.Token0Name, p.Token0Symbol, p.Token0Decimals,
			p.Token1Name, p.Token1Symbol, p.Token1Decimals,
			p.CreatedAt, p.UpdatedAt,
		}...)
	}

	sqlString := fmt.Sprintf(`
		insert into 
    		pools(address, token0, token1, tvl, trade_points, liquidity_points, 
    		      token0_name, token0_symbol, token0_decimals, 
    		      token1_name, token1_symbol, token1_decimals,
    		      created_at, updated_at)
		values
			%s
		on duplicate key update
			address = values(address),
			token0 = values(token0),
			token1 = values(token1),
			tvl = values(tvl),
			trade_points = trade_points + values(trade_points),
			liquidity_points = liquidity_points + values(liquidity_points),
			token0_name = values(token0_name),
			token0_symbol = values(token0_symbol),
			token0_decimals = values(token0_decimals),
			token1_name = values(token1_name),
			token1_symbol = values(token1_symbol),
			token1_decimals = values(token1_decimals),
			created_at = values(created_at),
			updated_at = values(updated_at)                      
	`, placeholders)

	return db.Exec(sqlString, params...).Error
}

func (service *PoolService) List(request model.PoolPagingRequest) (total int64, pools []*model.PoolInfo, err error) {
	var db *gorm.DB
	var sortField string

	if request.SortField == "tvl" {
		sortField = request.SortField
		db = service.store.DB.Model(&model.Pool{}).
			Select("pools.*, pool_params.trade_weight, pool_params.liquidity_weight").
			Joins("INNER JOIN pool_params ON pools.address = pool_params.address")
	} else {
		sortField = fmt.Sprintf("%s_weight", request.SortField)
		db = service.store.DB.Model(&model.PoolParams{}).
			Select("pool_params.*, pools.token0_symbol, pools.token1_symbol, pools.tvl").
			Joins("INNER JOIN pools ON pool_params.address = pools.address")
	}

	if err = db.Count(&total).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get count of pools")
	}

	var orderBy string
	if request.IsDesc() {
		orderBy = fmt.Sprintf("%s DESC", sortField)
	} else {
		orderBy = fmt.Sprintf("%s ASC", sortField)
	}

	if err = db.Order(orderBy).Offset(request.Offset).Limit(request.Limit).Find(&pools).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get pools")
	}

	return
}
