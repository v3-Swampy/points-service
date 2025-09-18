package service

import (
	"fmt"
	"time"

	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/sync"
	"gorm.io/gorm"
)

type PoolWeight struct {
	tradeWeight     int64 `default:"1"`
	liquidityWeight int64 `default:"1"`
}

type Config struct {
	PoolWeights map[string]PoolWeight
}

type StatService struct {
	config Config
	store  *store.Store
}

func NewStatService(config Config, store *store.Store) *StatService {
	return &StatService{
		config: config,
		store:  store,
	}
}

func (service *StatService) OnEventBatch(trades []sync.TradeEvent, liquidities []sync.LiquidityEvent) error {
	users := make(map[string]*model.User)
	pools := make(map[string]*model.Pool)

	service.aggregateTrade(trades, users, pools)
	service.aggregateLiquidity(liquidities, users, pools)

	if err := service.Store(users, pools); err != nil {
		return err
	}

	return nil
}

func (service *StatService) aggregateTrade(trades []sync.TradeEvent, users map[string]*model.User, pools map[string]*model.Pool) {
	for _, trade := range trades {
		statTime := time.Unix(trade.Timestamp, 0)
		user := trade.User
		pool := trade.Pool.TokenLP.Address.String()
		tradeWeight := service.config.PoolWeights[pool].tradeWeight
		tradePoints := trade.Value0.Add(trade.Value1).Mul(decimal.NewFromInt(tradeWeight))

		if u, exists := users[user]; exists {
			u.TradePoints = u.TradePoints.Add(tradePoints)
			u.UpdatedAt = statTime
		} else {
			users[user] = model.NewUser(user, tradePoints, decimal.Zero, statTime)
		}

		if p, exists := pools[pool]; exists {
			p.TradePoints = p.TradePoints.Add(tradePoints)
			p.UpdatedAt = statTime
		} else {
			pools[pool] = model.NewPool(trade.Pool, tradePoints, decimal.Zero, statTime)
		}
	}
}

func (service *StatService) aggregateLiquidity(liquidities []sync.LiquidityEvent, users map[string]*model.User, pools map[string]*model.Pool) {
	for _, liquidity := range liquidities {
		statTime := time.Unix(liquidity.Timestamp, 0)
		user := liquidity.User
		pool := liquidity.Pool.TokenLP.Address.String()
		liquidityWeight := service.config.PoolWeights[pool].liquidityWeight
		liquidityPoints := liquidity.ValueSecs.Mul(decimal.NewFromInt(liquidityWeight))

		if u, exists := users[user]; exists {
			u.LiquidityPoints = u.LiquidityPoints.Add(liquidityPoints)
			u.UpdatedAt = statTime
		} else {
			model.NewUser(user, decimal.Zero, liquidityPoints, statTime)
		}

		if p, exists := pools[pool]; exists {
			p.LiquidityPoints = p.LiquidityPoints.Add(liquidityPoints)
			p.UpdatedAt = statTime
		} else {
			pools[pool] = model.NewPool(liquidity.Pool, decimal.Zero, liquidityPoints, statTime)
		}
	}
}

func (service *StatService) Store(users map[string]*model.User, pools map[string]*model.Pool) error {
	service.store.DB.Transaction(func(dbTx *gorm.DB) error {
		if len(users) > 0 {
			userArray := make([]*model.User, 0, len(users))
			for _, user := range users {
				userArray = append(userArray, user)
			}
			if err := service.BatchDeltaUpsertUsers(service.store.DB, userArray); err != nil {
				return errors.WithMessage(err, "failed to batch delta upsert users")
			}
		}

		if len(pools) > 0 {
			poolArray := make([]*model.Pool, 0, len(pools))
			for _, pool := range pools {
				poolArray = append(poolArray, pool)
			}
			if err := service.BatchDeltaUpsertPools(service.store.DB, poolArray); err != nil {
				return errors.WithMessage(err, "failed to batch delta upsert pools")
			}
		}
		return nil
	})

	return nil
}

func (service *StatService) BatchDeltaUpsertUsers(dbTx *gorm.DB, users []*model.User) error {
	db := service.store.DB
	if dbTx != nil {
		db = dbTx
	}

	var placeholders string
	var params []interface{}
	size := len(users)
	for i, u := range users {
		placeholders += "(?,?,?,?,?)"
		if i != size-1 {
			placeholders += ",\n\t\t\t"
		}
		params = append(params, []interface{}{u.Address, u.TradePoints, u.LiquidityPoints, u.CreatedAt, u.UpdatedAt}...)
	}

	sqlString := fmt.Sprintf(`
		insert into 
    		users(address, trade_points, liquidity_points, created_at, updated_at)
		values
			%s
		on duplicate key update
			address = values(address),
			trade_points = trade_points + values(trade_points),
			liquidity_points = liquidity_points + values(liquidity_points),
			created_at = values(created_at),
			updated_at = values(updated_at)
	`, placeholders)

	if err := db.Exec(sqlString, params...).Error; err != nil {
		return err
	}

	return nil
}

func (service *StatService) BatchDeltaUpsertPools(dbTx *gorm.DB, pools []*model.Pool) error {
	db := service.store.DB
	if dbTx != nil {
		db = dbTx
	}

	var placeholders string
	var params []interface{}
	size := len(pools)
	for i, p := range pools {
		placeholders += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		if i != size-1 {
			placeholders += ",\n\t\t\t"
		}
		params = append(params, []interface{}{
			p.Address, p.Token0, p.Token1, p.Tvl, p.TradePoints, p.LiquidityPoints,
			p.TokenLpName, p.TokenLpSymbol, p.TokenLpDecimals, p.Token0Name, p.Token0Symbol, p.Token0Decimals,
			p.Token1Name, p.Token1Symbol, p.Token1Decimals, p.CreatedAt, p.UpdatedAt,
		}...)
	}

	sqlString := fmt.Sprintf(`
		insert into 
    		pools(address, token0, token1, tvl, trade_points, liquidity_points, 
    		      token_lp_name, token_lp_symbol, token_lp_decimals, token0_name, token0_symbol, token0_decimals, 
    		      token1_name, token1_symbol, token1_decimals, created_at, updated_at)
		values
			%s
		on duplicate key update
			address = values(address),
			token0 = values(token0),
			token1 = values(token1),
			tvl = values(tvl),
			trade_points = trade_points + values(trade_points),
			liquidity_points = liquidity_points + values(liquidity_points),                             
			token_lp_name = values(token_lp_name),
			token_lp_symbol = values(token_lp_symbol),
			token_lp_decimals = values(token_lp_decimals),
			token0_name = values(token0_name),
			token0_symbol = values(token0_symbol),
			token0_decimals = values(token0_decimals),
			token1_name = values(token1_name),
			token1_symbol = values(token1_symbol),
			token1_decimals = values(token1_decimals),
			created_at = values(created_at),
			updated_at = values(updated_at)
	`, placeholders)

	if err := db.Exec(sqlString, params...).Error; err != nil {
		return err
	}

	return nil
}
