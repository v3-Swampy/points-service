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

var logger = logrus.WithField("module", "parsing")

type Config struct {
	Endpoint          string
	NextHourTimestamp int64         // unix timestamp in seconds that truncated by hour
	PollInterval      time.Duration `default:"1m"`
}

type Service struct {
	*Client

	config  Config
	handler sync.EventHandler
	vswap   *blockchain.Vswap
	swappi  *blockchain.Swappi
	scan    *scan.Api
	pools   []common.Address
}

func NewService(config Config, handler sync.EventHandler, vswap *blockchain.Vswap, swappi *blockchain.Swappi, scan *scan.Api, pools ...common.Address) (*Service, error) {
	// TODO read from contract parser for the first hourTimestamp
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
		return nil, errors.WithMessagef(err, "Failed to create RPC client at %v", config.Endpoint)
	}

	return &Service{
		Client:  contractParsingClient,
		config:  config,
		handler: handler,
		vswap:   vswap,
		swappi:  swappi,
		scan:    scan,
		pools:   pools,
	}, nil
}

func (service *Service) Run(ctx context.Context, wg *stdSync.WaitGroup) {
	defer wg.Done()

	nextHourTimestamp := service.config.NextHourTimestamp
	logger.WithField("next", formatHourTimestamp(nextHourTimestamp)).Info("Start to poll data from contract parser")

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger := logger.WithField("dt", formatHourTimestamp(nextHourTimestamp)).WithField("ts", nextHourTimestamp)

			if ok, err := service.sync(ctx, nextHourTimestamp); err != nil {
				logger.WithError(err).Warn("Failed to sync once")
				ticker.Reset(service.config.PollInterval)
			} else if ok {
				nextHourTimestamp += 3600
				logger.Debug("Move forward")
				ticker.Reset(time.Millisecond)
			} else {
				logger.Debug("Data unavailable yet")
				ticker.Reset(service.config.PollInterval)
			}
		}
	}
}

func (service *Service) sync(ctx context.Context, hourTimestamp int64) (bool, error) {
	logger := logger.WithField("dt", formatHourTimestamp(hourTimestamp)).WithField("ts", hourTimestamp)

	// TODO cacheable
	minBlockNumber, err := service.scan.GetBlockNumberByTime(hourTimestamp, true)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to query min block number by hour timestamp from scan")
	}

	maxBlockNumber, err := service.scan.GetBlockNumberByTime(hourTimestamp+3600, false)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to query max block number by hour timestamp from scan")
	}

	logger.WithFields(logrus.Fields{
		"min": minBlockNumber,
		"max": maxBlockNumber,
	}).Debug("Block number range retrieved")

	var tradeEvents []sync.TradeEvent
	var liquidityEvents []sync.LiquidityEvent
	priceCache := make(map[common.Address]decimal.Decimal)

	for i, pool := range service.pools {
		logger.WithField("pool", pool).Debugf("Begin to collect data for pool [%v/%v]", i+1, len(service.pools))

		// retrieve trade data
		trades, err := service.GetHourlyTradeDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get trade data from contract parser")
		}

		if trades == nil {
			return false, nil
		}

		// retrieve liquidity data
		liquidities, err := service.GetHourlyLiquidityDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get liquidity data from contract parser")
		}

		if liquidities == nil {
			return false, nil
		}

		if len(trades) == 0 && len(liquidities) == 0 {
			continue
		}

		// get pool info
		info, err := service.vswap.GetPoolInfo(pool)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get pool info")
		}

		logger.WithField("pool", info).Debug("Pool info retrieved")

		// get prices to construct events
		price0, cached, err := service.getPrice(minBlockNumber, maxBlockNumber, pool, info.Token0.Address, priceCache)
		if err != nil {
			return false, errors.WithMessagef(err, "Failed to get price of token0 %v", info.Token0.Symbol)
		}

		if !cached {
			logger.WithField("price", price0.Truncate(6)).WithField("token", info.Token0.Symbol).Debug("Token0 price retrieved")
		}

		price1, cached, err := service.getPrice(minBlockNumber, maxBlockNumber, pool, info.Token1.Address, priceCache)
		if err != nil {
			return false, errors.WithMessagef(err, "Failed to get price of token1 %v", info.Token1.Symbol)
		}

		if !cached {
			logger.WithField("price", price1.Truncate(6)).WithField("token", info.Token1.Symbol).Debug("Token1 price retrieved")
		}

		// trade events
		for _, v := range trades {
			tradeEvents = append(tradeEvents, sync.TradeEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: hourTimestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				Value0: decimal.NewFromBigInt(v.Token0Volume.ToInt(), -int32(info.Token0.Decimals)).Mul(price0),
				Value1: decimal.NewFromBigInt(v.Token1Volume.ToInt(), -int32(info.Token1.Decimals)).Mul(price1),
			})
		}

		// liquidity events
		for _, v := range liquidities {
			liquidityEvents = append(liquidityEvents, sync.LiquidityEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: hourTimestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				Value0Secs: decimal.NewFromBigInt(v.Token0LiquiditySeconds.ToInt(), -int32(info.Token0.Decimals)).Mul(price0),
				Value1Secs: decimal.NewFromBigInt(v.Token1LiquiditySeconds.ToInt(), -int32(info.Token1.Decimals)).Mul(price1),
			})
		}
	}

	logger.WithFields(logrus.Fields{
		"trades":      len(tradeEvents),
		"liquidities": len(liquidityEvents),
	}).Debug("Trade and liquidity events retrieved")

	time := sync.TimeInfo{
		HourTimestamp:  hourTimestamp,
		MinBlockNumber: minBlockNumber,
		MaxBlockNumber: maxBlockNumber,
	}
	if err := service.handler.OnEventBatch(time, tradeEvents, liquidityEvents); err != nil {
		return false, errors.WithMessage(err, "Failed to handle trade and liquidity events")
	}

	logger.WithFields(logrus.Fields{
		"trades":      len(tradeEvents),
		"liquidities": len(liquidityEvents),
	}).Info("Trade and liquidity events handled")

	return true, nil
}

func (service *Service) getPrice(minBlockNumber, maxBlockNumber uint64, pool, token common.Address, cache map[common.Address]decimal.Decimal) (decimal.Decimal, bool, error) {
	if price, ok := cache[token]; ok {
		return price, true, nil
	}

	sumPrices := decimal.Zero
	step := (maxBlockNumber - minBlockNumber + 1) / 6
	var count int64

	for bn := minBlockNumber + step; bn < maxBlockNumber; bn += step {
		opts := bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(bn),
		}

		price, err := service.queryPrice(&opts, pool, token)
		if err != nil {
			return decimal.Zero, false, errors.WithMessagef(err, "Failed to sample token price at block %v", bn)
		}

		sumPrices = sumPrices.Add(price)
		count++
	}

	price := sumPrices.Div(decimal.NewFromInt(count))
	cache[token] = price

	return price, false, nil
}

func (service *Service) queryPrice(opts *bind.CallOpts, pool, token common.Address) (decimal.Decimal, error) {
	// get from swappi
	price, err := service.swappi.GetTokenPriceAuto(opts, token)
	if err == nil {
		return price, nil
	}

	if err != blockchain.ErrSwappiPairNotFound {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token price from Swappi")
	}

	// get from vswap
	price, err = service.vswap.GetTokenPriceUSDT(opts, pool, token)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token price from vSwap")
	}

	return price, nil
}

func formatHourTimestamp(hourTimestamp int64) string {
	return time.Unix(hourTimestamp, 0).Format(time.DateTime)
}
