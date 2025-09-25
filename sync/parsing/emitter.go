package parsing

import (
	"context"
	sdtErrors "errors"
	"math/big"
	stdSync "sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/blockchain"
	"github.com/v3-Swampy/points-service/sync"
)

var DefaultPriceSamples = uint64(4)

// Emitter is used to generate event based on polled data from contract parser.
type Emitter struct {
	buf    chan sync.BatchEvent
	vswap  *blockchain.Vswap
	swappi *blockchain.Swappi
	logger *logrus.Entry
}

func NewEmitter(vswap *blockchain.Vswap, swappi *blockchain.Swappi) *Emitter {
	return &Emitter{
		buf:    make(chan sync.BatchEvent, DefaultBufSize),
		vswap:  vswap,
		swappi: swappi,
		logger: logrus.WithField("worker", "sync.emitter"),
	}
}

func (emitter *Emitter) Close() {
	close(emitter.buf)
}

func (emitter *Emitter) Ch() <-chan sync.BatchEvent {
	return emitter.buf
}

func (emitter *Emitter) Run(ctx context.Context, wg *stdSync.WaitGroup, dataCh <-chan Snapshot) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-dataCh:
			emitter.mustEmit(ctx, data)
		}
	}
}

func (emitter *Emitter) mustEmit(ctx context.Context, data Snapshot) {
	logger := emitter.logger.WithFields(logrus.Fields{
		"ts":    formatTs(data.Timestamp),
		"minBN": data.MinBlockNumber,
		"maxBN": data.MaxBlockNumber,
	})

	for {
		start := time.Now()

		event, err := emitter.emit(ctx, data)
		if err != nil {
			logger.WithError(err).Warn("Failed to emit event")

			select {
			case <-ctx.Done():
				return
			case <-time.After(DefaultIntervalError):
				logger.Debug("Emitter retry to emit event")
			}
		} else {
			select {
			case emitter.buf <- event:
				logger.WithField("elapsed", time.Since(start)).Info("Emitter move forward")
				return
			case <-ctx.Done():
				return
			}
		}
	}
}

func (emitter *Emitter) emit(ctx context.Context, data Snapshot) (sync.BatchEvent, error) {
	logger := emitter.logger.WithField("ts", formatTs(data.Timestamp))

	event := sync.BatchEvent{
		TimeInfo: data.TimeInfo,
	}

	priceCache := make(map[common.Address]decimal.Decimal)
	interval := decimal.NewFromInt(data.IntervalSecs)

	for i, pool := range data.Pools {
		// check cancellation
		select {
		case <-ctx.Done():
			return sync.BatchEvent{}, ctx.Err()
		default:
		}

		logger.WithField("pool", pool.Address).Debugf("Begin to emit for pool [%v/%v]", i+1, len(data.Pools))

		if len(pool.Trades) == 0 && len(pool.Liquidities) == 0 {
			continue
		}

		// get pool info
		info, err := emitter.vswap.GetPoolInfo(pool.Address)
		if err != nil {
			return sync.BatchEvent{}, errors.WithMessage(err, "Failed to get pool info")
		}

		logger.WithField("pool", info).Debug("Pool info retrieved")

		// get prices to construct events
		price0, cached, err := emitter.getPrice(data.MinBlockNumber, data.MaxBlockNumber, pool.Address, info.Token0.Address, priceCache)
		if err != nil {
			return sync.BatchEvent{}, errors.WithMessagef(err, "Failed to get price of token0 %v", info.Token0.Symbol)
		}

		if !cached {
			logger.WithField("price", price0.Truncate(6)).WithField("token", info.Token0.Symbol).Debug("Token0 price retrieved")
		}

		price1, cached, err := emitter.getPrice(data.MinBlockNumber, data.MaxBlockNumber, pool.Address, info.Token1.Address, priceCache)
		if err != nil {
			return sync.BatchEvent{}, errors.WithMessagef(err, "Failed to get price of token1 %v", info.Token1.Symbol)
		}

		if !cached {
			logger.WithField("price", price1.Truncate(6)).WithField("token", info.Token1.Symbol).Debug("Token1 price retrieved")
		}

		// trade events
		for _, v := range pool.Trades {
			event.Trades = append(event.Trades, sync.TradeEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: data.Timestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				Value0: decimal.NewFromBigInt(v.Token0Volume.ToInt(), -int32(info.Token0.Decimals)).Mul(price0),
				Value1: decimal.NewFromBigInt(v.Token1Volume.ToInt(), -int32(info.Token1.Decimals)).Mul(price1),
			})
		}

		// liquidity events
		for _, v := range pool.Liquidities {
			event.Liquidities = append(event.Liquidities, sync.LiquidityEvent{
				PoolEvent: sync.PoolEvent{
					Timestamp: data.Timestamp,
					User:      v.UserAddress,
					Pool:      info,
				},
				Value0: decimal.NewFromBigInt(v.Token0LiquiditySeconds.ToInt(), -int32(info.Token0.Decimals)).Mul(price0).Div(interval),
				Value1: decimal.NewFromBigInt(v.Token1LiquiditySeconds.ToInt(), -int32(info.Token1.Decimals)).Mul(price1).Div(interval),
			})
		}
	}

	logger.WithFields(logrus.Fields{
		"trades":      len(event.Trades),
		"liquidities": len(event.Liquidities),
	}).Debug("Trade and liquidity events generated")

	return event, nil
}

func (emitter *Emitter) getPrice(minBlockNumber, maxBlockNumber uint64, pool, token common.Address, cache map[common.Address]decimal.Decimal) (decimal.Decimal, bool, error) {
	if price, ok := cache[token]; ok {
		return price, true, nil
	}

	if minBlockNumber > maxBlockNumber {
		emitter.logger.WithFields(logrus.Fields{
			"min": minBlockNumber,
			"max": maxBlockNumber,
		}).Error("Invalid block number range")
	}

	// sample n prices
	step := (maxBlockNumber + 1 - minBlockNumber) / DefaultPriceSamples
	if step == 0 {
		step = 1
	}

	sumPrices := decimal.Zero
	var count int64

	// ensure the maxBlockNumber sampled in case that liquidity added at maxBlockNumber
	for bn := maxBlockNumber; bn >= minBlockNumber && bn <= maxBlockNumber; bn -= step {
		opts := bind.CallOpts{
			BlockNumber: new(big.Int).SetUint64(bn),
		}

		price, err := emitter.queryPrice(&opts, pool, token)
		if err != nil {
			return decimal.Zero, false, errors.WithMessagef(err, "Failed to sample token price at block %v", bn)
		}

		if price.IsZero() {
			break
		}

		sumPrices = sumPrices.Add(price)
		count++
	}

	if count == 0 {
		emitter.logger.WithFields(logrus.Fields{
			"minBN": minBlockNumber,
			"maxBN": maxBlockNumber,
		}).Fatal("No token price sampled")
	}

	price := sumPrices.Div(decimal.NewFromInt(count))
	cache[token] = price

	return price, false, nil
}

func (emitter *Emitter) queryPrice(opts *bind.CallOpts, pool, token common.Address) (decimal.Decimal, error) {
	// get from swappi
	price, err := emitter.swappi.GetTokenPriceAuto(opts, token)
	if err == nil {
		return price, nil
	}

	if !sdtErrors.Is(err, blockchain.ErrSwappiPairNotFound) {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token price from Swappi")
	}

	// get from vswap
	price, err = emitter.vswap.GetTokenPriceUSDT(opts, pool, token)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token price from vSwap")
	}

	return price, nil
}
