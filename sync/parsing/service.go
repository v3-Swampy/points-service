package parsing

import (
	"context"
	stdSync "sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/sync"
	"github.com/v3-Swampy/points-service/sync/blockchain"
)

type Config struct {
	Endpoint          string
	NextHourTimestamp int64         // unix timestamp in seconds that truncated by hour
	PollInterval      time.Duration `default:"1m"`
	Pools             []string
	Swappi            blockchain.SwappiAddresses
}

type Service struct {
	*Client

	config Config

	handler sync.EventHandler

	swappi *blockchain.Swappi
}

func NewService(config Config, handler sync.EventHandler, client *web3go.Client) (*Service, error) {
	if config.NextHourTimestamp%3600 > 0 {
		return nil, errors.Errorf("Invalid NextHourTimestamp value %v", config.NextHourTimestamp)
	}

	if len(config.Pools) == 0 {
		return nil, errors.New("Pools not specified")
	}

	for _, v := range config.Pools {
		if !common.IsHexAddress(v) {
			return nil, errors.Errorf("Invalid hex address %v", v)
		}
	}

	contractParsingClient, err := NewClient(config.Endpoint)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to create RPC client")
	}

	caller, _ := client.ToClientForContract()
	erc20 := blockchain.NewERC20(caller)
	swappi := blockchain.NewSwappi(caller, erc20, config.Swappi)

	return &Service{
		Client:  contractParsingClient,
		config:  config,
		handler: handler,
		swappi:  swappi,
	}, nil
}

func (service *Service) Run(ctx context.Context, wg *stdSync.WaitGroup) {
	defer wg.Done()

	// TODO load `next` value from DB or default configured
	nextHourTimestamp := service.config.NextHourTimestamp

	ticker := time.NewTicker(service.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ok, err := service.sync(ctx, nextHourTimestamp); err != nil {
				logrus.WithError(err).WithField("hourTimestamp", nextHourTimestamp).Warn("Failed to sync once")
			} else if ok {
				nextHourTimestamp += 3600
			}
		}
	}
}

func (service *Service) sync(ctx context.Context, hourTimestamp int64) (bool, error) {
	var tradeEvents []sync.TradeEvent
	var liquidityEvents []sync.LiquidityEvent
	priceCache := make(map[common.Address]decimal.Decimal)

	for _, pool := range service.config.Pools {
		// retrieve trade data
		trades, err := service.GetHourlyTradeDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get trade data")
		}

		if trades == nil {
			return false, nil
		}

		// retrieve liquidity data
		liquidities, err := service.GetHourlyLiquidityDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get liquidity data")
		}

		if liquidities == nil {
			return false, nil
		}

		// construct events to handle
		info, err := service.swappi.GetPairInfo(common.HexToAddress(pool))
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get pool info")
		}

		// TODO sample and get average prices in the given hour time
		price0, err := service.getPrice(nil, info.Token0.Address, false, priceCache)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get price of token0")
		}

		price1, err := service.getPrice(nil, info.Token1.Address, false, priceCache)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get price of token1")
		}

		// trade events
		for _, v := range trades {
			tradeEvents = append(tradeEvents, sync.TradeEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: hourTimestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				Value0: decimal.NewFromBigInt(v.Token0Volumes.ToInt(), -int32(info.Token0.Decimals)).Mul(price0),
				Value1: decimal.NewFromBigInt(v.Token1Volumes.ToInt(), -int32(info.Token1.Decimals)).Mul(price1),
			})
		}

		// priceLP, err := service.getPrice(nil, info.TokenLP.Address, true, priceCache)
		// if err != nil {
		// 	return false, errors.WithMessage(err, "Failed to get price of LP token")
		// }

		// liquidity events
		for _, v := range liquidities {
			liquidityEvents = append(liquidityEvents, sync.LiquidityEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: hourTimestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				// ValueSecs: decimal.NewFromBigInt(v.LiquiditySeconds.ToInt(), -int32(info.TokenLP.Decimals)).Mul(priceLP),
				Value0Secs: decimal.NewFromBigInt(v.Token0LiquiditySeconds.ToInt(), -int32(info.Token0.Decimals)).Mul(price0),
				Value1Secs: decimal.NewFromBigInt(v.Token1LiquiditySeconds.ToInt(), -int32(info.Token1.Decimals)).Mul(price1),
			})
		}
	}

	if err := service.handler.OnEventBatch(hourTimestamp, tradeEvents, liquidityEvents); err != nil {
		return false, errors.WithMessage(err, "Failed to handle trade and liquidity events")
	}

	logrus.WithFields(logrus.Fields{
		"timestamp":   time.Unix(hourTimestamp, 0).Format(time.DateTime),
		"trades":      len(tradeEvents),
		"liquidities": len(liquidityEvents),
	}).Info("Succeeded to handle trade and liquidity events")

	return true, nil
}

func (service *Service) getPrice(opts *bind.CallOpts, token common.Address, isLP bool, cache map[common.Address]decimal.Decimal) (price decimal.Decimal, err error) {
	if price, ok := cache[token]; ok {
		return price, nil
	}

	if isLP {
		price, err = service.swappi.GetPairTokenPrice(opts, token)
	} else {
		price, err = service.swappi.GetTokenPriceAuto(opts, token)
	}

	if err == nil {
		cache[token] = price
	}

	return
}
