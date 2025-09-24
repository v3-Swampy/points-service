package parsing

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/v3-Swampy/points-service/blockchain/scan"
)

var (
	DefaultBufSize       = 1024
	DefaultIntervalError = 5 * time.Second
	DefaultIntervalIdle  = time.Minute
)

type PollConfig struct {
	Endpoint string // RPC endpoint of contract parser
	Scan     string // Open API endpoint of Scan
}

// Poller is used to poll trade and liquidity data from contract parser.
type Poller struct {
	client            *Client
	scan              *scan.Api
	buf               chan HourlyData
	nextHourTimestamp int64
	pools             []common.Address
}

// NewPoller creates a new poller.
//
// If the given nextHourTimestamp is 0, then retrieve the first timestamp from contract parser.
//
// Note, it returns error if the given pools is empty.
func NewPoller(config PollConfig, nextHourTimestamp int64, pools ...common.Address) (*Poller, error) {
	if len(pools) == 0 {
		return nil, errors.New("Pools not specified")
	}

	client, err := NewClient(config.Endpoint)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create client")
	}

	// retrieve first timestamp
	if nextHourTimestamp == 0 {
		if nextHourTimestamp, err = client.FirstTimestamp(context.Background()); err != nil {
			return nil, errors.WithMessage(err, "Failed to poll first timestamp")
		}
	}

	return &Poller{
		client:            client,
		scan:              scan.NewApi(config.Scan),
		buf:               make(chan HourlyData, DefaultBufSize),
		nextHourTimestamp: nextHourTimestamp,
		pools:             pools,
	}, nil
}

func (poller *Poller) Close() {
	poller.client.Close()
	close(poller.buf)
}

func (poller *Poller) Ch() <-chan HourlyData {
	return poller.buf
}

func (poller *Poller) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	hourTimestamp := poller.nextHourTimestamp
	logHourTimestamp(hourTimestamp).Info("Poller started")

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger = logHourTimestamp(hourTimestamp)

			start := time.Now()

			data, ok, err := poller.poll(ctx, hourTimestamp)
			if err != nil {
				logger.WithError(err).Warn("Failed to poll data from contract parser")
				ticker.Reset(DefaultIntervalError)
			} else if ok {
				select {
				case poller.buf <- data:
					logger.WithField("elapsed", time.Since(start)).Info("Poller move forward")
					hourTimestamp += 3600
				case <-ctx.Done():
					return
				}
				ticker.Reset(time.Millisecond)
			} else {
				logger.Debug("Poller is idle")
				ticker.Reset(DefaultIntervalIdle)
			}
		}
	}
}

func (poller *Poller) poll(ctx context.Context, hourTimestamp int64) (HourlyData, bool, error) {
	var result HourlyData
	result.HourTimestamp = hourTimestamp

	for _, pool := range poller.pools {
		// trade data
		trades, err := poller.client.GetHourlyTradeDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return HourlyData{}, false, errors.WithMessage(err, "Failed to poll trade data")
		}

		if trades == nil {
			return HourlyData{}, false, nil
		}

		// liquidity data
		liquidities, err := poller.client.GetHourlyLiquidityDataAll(ctx, pool, hourTimestamp)
		if err != nil {
			return HourlyData{}, false, errors.WithMessage(err, "Failed to poll liquidity data")
		}

		if liquidities == nil {
			return HourlyData{}, false, nil
		}

		result.Pools = append(result.Pools, PoolData{
			Address:     pool,
			Trades:      trades,
			Liquidities: liquidities,
		})
	}

	var err error

	if result.MinBlockNumber, err = poller.scan.GetBlockNumberByTime(hourTimestamp, true); err != nil {
		return HourlyData{}, false, errors.WithMessage(err, "Failed to query min block number")
	}

	if result.MaxBlockNumber, err = poller.scan.GetBlockNumberByTime(hourTimestamp+3600, false); err != nil {
		return HourlyData{}, false, errors.WithMessage(err, "Failed to query max block number")
	}

	return result, true, nil
}
