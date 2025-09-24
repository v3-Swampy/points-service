package parsing

import (
	"context"
	stdSync "sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/sync"
)

const (
	DefaultEventBatchSize    = 100
	DefaultEventBatchTimeout = time.Second * 10
)

type Batcher struct {
	handler sync.EventHandler
}

func NewBatcher(handler sync.EventHandler) *Batcher {
	return &Batcher{handler}
}

func (batcher *Batcher) Run(ctx context.Context, wg *stdSync.WaitGroup, eventCh <-chan sync.BatchEvent) {
	defer wg.Done()

	var batch sync.BatchEvent

	ticker := time.NewTicker(DefaultEventBatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			batch = batcher.mustHandle(ctx, batch)
			ticker.Reset(DefaultEventBatchTimeout)
		case event := <-eventCh:
			batch.Merge(event)

			if len(batch.Trades)+len(batch.Liquidities) >= DefaultEventBatchSize {
				batch = batcher.mustHandle(ctx, batch)
				ticker.Reset(DefaultEventBatchTimeout)
			}
		}
	}
}
func (batcher *Batcher) mustHandle(ctx context.Context, batch sync.BatchEvent) sync.BatchEvent {
	if batch.HourTimestamp == 0 {
		return sync.BatchEvent{}
	}

	logger = logger.WithFields(logrus.Fields{
		"ts": batch.HourTimestamp,
		"dt": formatHourTimestamp(batch.HourTimestamp),
	})

	for {
		start := time.Now()

		if err := batcher.handler.OnEventBatch(batch); err != nil {
			logger.WithError(err).Warn("Failed to handle events in batch")

			select {
			case <-ctx.Done():
				return sync.BatchEvent{}
			case <-time.After(DefaultIntervalError):
				logger.Debug("Batcher retry to handle events")
			}
		} else {
			logger.WithFields(logrus.Fields{
				"elapsed":   time.Since(start),
				"trade":     len(batch.Trades),
				"liquidity": len(batch.Liquidities),
			}).Info("Batch events handled")

			return sync.BatchEvent{}
		}
	}
}
