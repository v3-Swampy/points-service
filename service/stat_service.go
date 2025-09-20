package service

import (
	"time"

	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/sync"
	"gorm.io/gorm"
)

type StatService struct {
	store  *store.Store
	config *ConfigService
	param  *PoolParamService
	user   *UserService
	pool   *PoolService
}

func NewStatService(store *store.Store) *StatService {
	return &StatService{
		store:  store,
		config: NewConfigService(store),
		param:  NewPoolParamService(store),
		user:   NewUserService(store),
		pool:   NewPoolService(store),
	}
}

func (service *StatService) OnEventBatch(timestamp int64, trades []sync.TradeEvent, liquidities []sync.LiquidityEvent) error {
	users := make(map[string]*model.User)
	pools := make(map[string]*model.Pool)

	service.aggregateTrade(trades, users, pools)
	service.aggregateLiquidity(liquidities, users, pools)

	return service.Store(timestamp, users, pools)
}

func (service *StatService) aggregateTrade(event []sync.TradeEvent, users map[string]*model.User, pools map[string]*model.Pool) error {
	for _, trade := range event {
		statTime := time.Unix(trade.Timestamp, 0)
		user := trade.User
		pool := trade.Pool.Address.String()

		weight, err := service.param.Get(pool)
		if err != nil {
			return err
		}
		tradeWeight := weight.TradeWeight
		tradePoints := trade.Value0.Add(trade.Value1).Mul(decimal.NewFromInt(int64(tradeWeight)))

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
	return nil
}

func (service *StatService) aggregateLiquidity(event []sync.LiquidityEvent, users map[string]*model.User, pools map[string]*model.Pool) error {
	for _, liquidity := range event {
		statTime := time.Unix(liquidity.Timestamp, 0)
		user := liquidity.User
		pool := liquidity.Pool.Address.String()

		weight, err := service.param.Get(pool)
		if err != nil {
			return err
		}
		liquidityWeight := weight.LiquidityWeight
		liquidityPoints := liquidity.Value0Secs.Add(liquidity.Value1Secs).Mul(decimal.NewFromInt(int64(liquidityWeight)))

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
	return nil
}

func (service *StatService) Store(timestamp int64, users map[string]*model.User, pools map[string]*model.Pool) error {
	service.store.DB.Transaction(func(dbTx *gorm.DB) error {
		if len(users) > 0 {
			userArray := make([]*model.User, 0, len(users))
			for _, user := range users {
				userArray = append(userArray, user)
			}
			if err := service.user.BatchDeltaUpsert(userArray, dbTx); err != nil {
				return errors.WithMessage(err, "failed to batch delta upsert users")
			}
		}

		if len(pools) > 0 {
			poolArray := make([]*model.Pool, 0, len(pools))
			for _, pool := range pools {
				poolArray = append(poolArray, pool)
			}
			if err := service.pool.BatchDeltaUpsert(poolArray, dbTx); err != nil {
				return errors.WithMessage(err, "failed to batch delta upsert pools")
			}
		}

		updateTime := time.Unix(timestamp, 0).Format(time.RFC3339)

		if err := service.config.StoreConfig(CfgKeyLastStatTimePoints, updateTime, dbTx); err != nil {
			return err
		}

		return nil
	})

	return nil
}
