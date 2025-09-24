package parsing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/blockchain/scan"
	"golang.org/x/sync/errgroup"
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
	logger            *logrus.Entry
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
		logger:            logrus.WithField("worker", "sync.poller"),
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
	poller.logger.WithField("ts", formatTs(hourTimestamp)).Info("Poller started")

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger := poller.logger.WithField("ts", formatTs(hourTimestamp))

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
	// check if data avaialbe
	latestHourTimestamp, err := poller.client.LatestTimestamp(ctx)
	if err != nil {
		return HourlyData{}, false, errors.WithMessage(err, "Failed to poll latest timestamp")
	}

	if hourTimestamp > latestHourTimestamp {
		return HourlyData{}, false, nil
	}

	// poll data in async
	var result HourlyData
	result.HourTimestamp = hourTimestamp

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)

	numPools := len(poller.pools)
	poolDataCh := make(chan PoolData, numPools)
	defer close(poolDataCh)

	// poll pool data
	for i := 0; i < numPools; i++ {
		pool := poller.pools[i]

		group.Go(func() (err error) {
			data := PoolData{
				Address: pool,
			}

			if data.Trades, err = poller.client.GetHourlyTradeDataAll(ctx, pool, hourTimestamp); err != nil {
				return errors.WithMessage(err, "Failed to poll trade data")
			}

			if data.Liquidities, err = poller.client.GetHourlyLiquidityDataAll(ctx, pool, hourTimestamp); err != nil {
				return errors.WithMessage(err, "Failed to poll liquidity data")
			}

			poolDataCh <- data

			return nil
		})
	}

	// poll min block number from scan
	group.Go(func() (err error) {
		if result.MinBlockNumber, err = poller.scan.GetBlockNumberByTime(hourTimestamp, true); err != nil {
			return errors.WithMessage(err, "Failed to query min block number")
		}

		return nil
	})

	// poll max block number from scan
	group.Go(func() (err error) {
		if result.MaxBlockNumber, err = poller.scan.GetBlockNumberByTime(hourTimestamp+3600, false); err != nil {
			return errors.WithMessage(err, "Failed to query max block number")
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		return HourlyData{}, false, errors.WithMessage(err, "Any async worker failed")
	}

	for i := 0; i < numPools; i++ {
		data := <-poolDataCh
		result.Pools = append(result.Pools, data)
	}

	return result, true, nil
}

func formatTs(hourTimestamp int64) string {
	dt := time.Unix(hourTimestamp, 0).Format(time.DateTime)
	return fmt.Sprintf("%v (%v)", dt, hourTimestamp)
}
