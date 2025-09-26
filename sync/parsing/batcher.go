package parsing

import (
	"context"
	stdSync "sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/v3-Swampy/points-service/sync"
)

type BatchOption struct {
	BatchSize     int           `default:"100"`
	BatchTimeout  time.Duration `default:"3s"`
	IntervalError time.Duration `default:"5s"`
}

type Batcher struct {
	option  BatchOption
	handler sync.EventHandler
	logger  *logrus.Entry
}

func NewBatcher(handler sync.EventHandler, option ...BatchOption) *Batcher {
	return &Batcher{
		option:  optionWithDefault(option...),
		handler: handler,
		logger:  logrus.WithField("worker", "sync.batcher"),
	}
}

func (batcher *Batcher) Run(ctx context.Context, wg *stdSync.WaitGroup, eventCh <-chan sync.BatchEvent) {
	defer wg.Done()

	var batch sync.BatchEvent

	ticker := time.NewTicker(batcher.option.BatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			batch = batcher.mustHandle(ctx, batch)
			ticker.Reset(batcher.option.BatchTimeout)
		case event := <-eventCh:
			batch.Merge(event)

			if len(batch.Trades)+len(batch.Liquidities) >= batcher.option.BatchSize {
				batch = batcher.mustHandle(ctx, batch)
				ticker.Reset(batcher.option.BatchTimeout)
			}
		}
	}
}
func (batcher *Batcher) mustHandle(ctx context.Context, batch sync.BatchEvent) sync.BatchEvent {
	if batch.Timestamp == 0 {
		return sync.BatchEvent{}
	}

	logger := batcher.logger.WithField("ts", formatTs(batch.Timestamp))

	for {
		start := time.Now()

		if err := batcher.handler.OnEventBatch(batch); err != nil {
			logger.WithError(err).Warn("Failed to handle events in batch")

			select {
			case <-ctx.Done():
				return sync.BatchEvent{}
			case <-time.After(batcher.option.IntervalError):
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
