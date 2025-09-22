package parsing

import (
	"context"
	"math/big"
	stdSync "sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/blockchain"
	"github.com/v3-Swampy/points-service/blockchain/scan"
	"github.com/v3-Swampy/points-service/sync"
)

type Config struct {
	Endpoint          string
	NextHourTimestamp int64         // unix timestamp in seconds that truncated by hour
	PollInterval      time.Duration `default:"1m"`
}

type Service struct {
	*Client

	config  Config
	handler sync.EventHandler
	swappi  *blockchain.Swappi
	scan    *scan.Api
	pools   []common.Address
}

func NewService(config Config, handler sync.EventHandler, swappi *blockchain.Swappi, scan *scan.Api, pools ...common.Address) (*Service, error) {
	if config.NextHourTimestamp == 0 {
		config.NextHourTimestamp = time.Now().Truncate(time.Hour).Unix()
	}

	if config.NextHourTimestamp%3600 > 0 {
		return nil, errors.Errorf("Invalid NextHourTimestamp value %v", config.NextHourTimestamp)
	}

	if len(pools) == 0 {
		return nil, errors.New("Pools not specified")
	}

	contractParsingClient, err := NewClient(config.Endpoint)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to create RPC client")
	}

	return &Service{
		Client:  contractParsingClient,
		config:  config,
		handler: handler,
		swappi:  swappi,
		scan:    scan,
		pools:   pools,
	}, nil
}

func (service *Service) Run(ctx context.Context, wg *stdSync.WaitGroup) {
	defer wg.Done()

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
	minBlockNumber, err := service.scan.GetBlockNumberByTime(hourTimestamp, true)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to query min block number by hour timestamp from scan")
	}

	maxBlockNumber, err := service.scan.GetBlockNumberByTime(hourTimestamp+3600, false)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to query max block number by hour timestamp from scan")
	}

	var tradeEvents []sync.TradeEvent
	var liquidityEvents []sync.LiquidityEvent
	priceCache := make(map[common.Address]decimal.Decimal)

	for _, pool := range service.pools {
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
		info, err := service.swappi.GetPairInfo(pool)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get pool info")
		}

		price0, err := service.getPrice(minBlockNumber, maxBlockNumber, info.Token0.Address, priceCache)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get price of token0")
		}

		price1, err := service.getPrice(minBlockNumber, maxBlockNumber, info.Token1.Address, priceCache)
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

func (service *Service) getPrice(minBlockNumber, maxBlockNumber uint64, token common.Address, cache map[common.Address]decimal.Decimal) (decimal.Decimal, error) {
	if price, ok := cache[token]; ok {
		return price, nil
	}

	sumPrices := decimal.Zero
	step := (maxBlockNumber - minBlockNumber + 1) / 6
	var count int64

	for bn := minBlockNumber + step; bn < maxBlockNumber; bn++ {
		opts := bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(bn),
		}

		price, err := service.swappi.GetTokenPriceAuto(&opts, token)
		if err != nil {
			return decimal.Zero, errors.WithMessage(err, "Failed to sample token price")
		}

		sumPrices = sumPrices.Add(price)
		count++
	}

	price := sumPrices.Div(decimal.NewFromInt(count))
	cache[token] = price

	return price, nil
}
