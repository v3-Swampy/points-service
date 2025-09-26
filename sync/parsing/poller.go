package parsing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/blockchain/scan"
	"golang.org/x/sync/errgroup"
)

type PollOption struct {
	BufferSize    int           `default:"1024"`
	IntervalError time.Duration `default:"5s"`
	IntervalIdle  time.Duration `default:"3s"`

	RPC  providers.Option
	Scan scan.Option
}

// Poller is used to poll trade and liquidity data from contract parser.
type Poller struct {
	option        PollOption
	client        *Client
	scan          *scan.Api
	buf           chan Snapshot
	nextTimestamp int64
	intervalSecs  int64
	pools         []common.Address
	logger        *logrus.Entry
}

// NewPoller creates a new poller.
//
// If the given lastTimestamp is 0, then retrieve the first timestamp from contract parser.
//
// Note, it returns error if the given pools is empty.
func NewPoller(rpcUrl, scanUrl string, lastTimestamp int64, pools []common.Address, option ...PollOption) (*Poller, error) {
	if len(pools) == 0 {
		return nil, errors.New("Pools not specified")
	}

	// init rpc client
	client, err := NewClient(rpcUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create client")
	}

	intervalSecs, err := client.SnapshotIntervalSecs()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get snapshot interval")
	}

	// retrieve first timestamp
	var nextTimestamp int64
	if lastTimestamp == 0 {
		if nextTimestamp, err = client.FirstTimestamp(); err != nil {
			return nil, errors.WithMessage(err, "Failed to poll first timestamp")
		}
	} else {
		nextTimestamp = lastTimestamp + intervalSecs
	}

	opt := optionWithDefault(option...)

	return &Poller{
		option:        opt,
		client:        client,
		scan:          scan.NewApi(scanUrl, opt.Scan),
		buf:           make(chan Snapshot, opt.BufferSize),
		nextTimestamp: nextTimestamp,
		intervalSecs:  intervalSecs,
		pools:         pools,
		logger:        logrus.WithField("worker", "sync.poller"),
	}, nil
}

func (poller *Poller) Close() {
	poller.client.Close()
	close(poller.buf)
}

func (poller *Poller) Ch() <-chan Snapshot {
	return poller.buf
}

func (poller *Poller) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	timestamp := poller.nextTimestamp
	poller.logger.WithField("ts", formatTs(timestamp)).Info("Poller started")

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	var lastMaxBlockNumber uint64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger := poller.logger.WithField("ts", formatTs(timestamp))

			start := time.Now()

			data, ok, err := poller.poll(timestamp, lastMaxBlockNumber)
			if err != nil {
				logger.WithError(err).Warn("Failed to poll data from contract parser")
				ticker.Reset(poller.option.IntervalError)
			} else if ok {
				select {
				case poller.buf <- data:
					logger.WithField("elapsed", time.Since(start)).Info("Poller move forward")
					timestamp += poller.intervalSecs
					lastMaxBlockNumber = data.MaxBlockNumber
				case <-ctx.Done():
					return
				}
				ticker.Reset(time.Millisecond)
			} else {
				logger.Debug("Poller is idle")
				ticker.Reset(poller.option.IntervalIdle)
			}
		}
	}
}

func (poller *Poller) poll(timestamp int64, lastMaxBlockNumber uint64) (Snapshot, bool, error) {
	// check if data avaialbe
	latestTimestamp, err := poller.client.LatestTimestamp()
	if err != nil {
		return Snapshot{}, false, errors.WithMessage(err, "Failed to poll latest timestamp")
	}

	if timestamp > latestTimestamp {
		return Snapshot{}, false, nil
	}

	// poll data in async
	var result Snapshot
	result.Timestamp = timestamp

	group := new(errgroup.Group)

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

			if data.Trades, err = poller.client.GetTradeDataAll(pool, timestamp); err != nil {
				return errors.WithMessage(err, "Failed to poll trade data")
			}

			if data.Liquidities, err = poller.client.GetLiquidityDataAll(pool, timestamp); err != nil {
				return errors.WithMessage(err, "Failed to poll liquidity data")
			}

			poolDataCh <- data

			return nil
		})
	}

	// poll min block number from scan
	var minBlockNumber uint64
	group.Go(func() error {
		if lastMaxBlockNumber > 0 {
			minBlockNumber = lastMaxBlockNumber + 1
			return nil
		}

		var startTime int64
		if timestamp > poller.intervalSecs {
			startTime = timestamp - poller.intervalSecs
		}

		bn, err := poller.scan.GetBlockNumberByTime(startTime, true)
		if err != nil {
			return errors.WithMessage(err, "Failed to query min block number")
		}

		if bn == 0 {
			return errors.Errorf("Failed to get min block number from scan, 0 returned by timestamp %v", startTime)
		}

		minBlockNumber = bn

		return nil
	})

	// poll max block number from scan
	var maxBlockNumber uint64
	group.Go(func() error {
		bn, err := poller.scan.GetBlockNumberByTime(timestamp-1, false)
		if err != nil {
			return errors.WithMessage(err, "Failed to query max block number")
		}

		if bn == 0 {
			return errors.Errorf("Failed to get max block number from scan, 0 returned by timestamp %v", timestamp)
		}

		maxBlockNumber = bn

		return nil
	})

	if err := group.Wait(); err != nil {
		return Snapshot{}, false, errors.WithMessage(err, "Any async worker failed")
	}

	result.MinBlockNumber = minBlockNumber
	result.MaxBlockNumber = maxBlockNumber

	for i := 0; i < numPools; i++ {
		data := <-poolDataCh
		result.Pools = append(result.Pools, data)
	}

	if result.MinBlockNumber == 0 || result.MaxBlockNumber == 0 || result.MinBlockNumber > result.MaxBlockNumber {
		poller.logger.WithFields(logrus.Fields{
			"min": result.MinBlockNumber,
			"max": result.MaxBlockNumber,
		}).Fatal("Invalid block number retrieved from scan")
	}

	return result, true, nil
}

func formatTs(timestamp int64) string {
	dt := time.Unix(timestamp, 0).Format(time.DateTime)
	return fmt.Sprintf("%v (%v)", dt, timestamp)
}
