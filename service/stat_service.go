package service

import (
	"time"

	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/model"
	"github.com/v3-Swampy/points-service/sync"
	"github.com/v3-Swampy/points-service/sync/blockchain"
	"gorm.io/gorm"
)

type SwappiConfig struct {
	Factory string
	USDT    string
	WCFX    string
}

type StatService struct {
	store  *store.Store
	config *ConfigService
	param  *PoolParamService
	user   *UserService
	pool   *PoolService
	swappi *blockchain.Swappi
}

func NewStatService(config SwappiConfig, client *web3go.Client, store *store.Store) *StatService {
	caller, _ := client.ToClientForContract()
	erc20 := blockchain.NewERC20(caller)
	swappi := blockchain.NewSwappi(caller, erc20, blockchain.SwappiAddresses{
		Factory: common.HexToAddress(config.Factory),
		USDT:    common.HexToAddress(config.USDT),
		WCFX:    common.HexToAddress(config.WCFX),
	})

	return &StatService{
		store:  store,
		config: NewConfigService(store),
		param:  NewPoolParamService(store),
		user:   NewUserService(store),
		pool:   NewPoolService(store),
		swappi: swappi,
	}
}

func (service *StatService) OnEventBatch(timestamp int64, trades []sync.TradeEvent, liquidities []sync.LiquidityEvent) error {
	users := make(map[string]*model.User)
	pools := make(map[string]*model.Pool)

	if err := service.aggregateTrade(trades, users, pools); err != nil {
		return err
	}

	if err := service.aggregateLiquidity(liquidities, users, pools); err != nil {
		return err
	}

	if err := service.UpdateTVL(pools); err != nil {
		return err
	}

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
		tradePoints := trade.Value0.Add(trade.Value1).Mul(decimal.NewFromInt(int64(weight.TradeWeight)))

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
		liquidityPoints := liquidity.Value0Secs.Add(liquidity.Value1Secs).Mul(decimal.NewFromFloat(0.1)).Mul(decimal.NewFromInt(int64(weight.LiquidityWeight)))

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

func (service *StatService) UpdateTVL(pools map[string]*model.Pool) error {
	if len(pools) > 0 {
		for _, pool := range pools {
			tvl, err := service.swappi.GetPairTVL(nil, common.HexToAddress(pool.Address))
			if err != nil {
				return err
			}
			pool.Tvl = tvl
		}
	}
	return nil
}

func (service *StatService) Store(timestamp int64, users map[string]*model.User, pools map[string]*model.Pool) error {
	return service.store.DB.Transaction(func(dbTx *gorm.DB) error {
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

		if err := service.config.UpsertLastStatPointsTime(timestamp, dbTx); err != nil {
			return err
		}

		return nil
	})
}
