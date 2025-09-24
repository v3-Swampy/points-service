package sync

import (
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/blockchain"
)

type PoolEvent struct {
	Timestamp int64  // unix timestamp in seconds
	User      string // user address, e.g. trader or liquidity provider
	Pool      blockchain.PairInfo
}

type TradeEvent struct {
	PoolEvent

	Value0 decimal.Decimal // amount0 * price0
	Value1 decimal.Decimal // amount1 * price1
}

type LiquidityEvent struct {
	PoolEvent

	Value0Secs decimal.Decimal // liquidity0 * price * secs
	Value1Secs decimal.Decimal // liquidity0 * price * secs
}

type TimeInfo struct {
	HourTimestamp  int64
	MinBlockNumber uint64
	MaxBlockNumber uint64
}

type BatchEvent struct {
	TimeInfo

	Trades      []TradeEvent
	Liquidities []LiquidityEvent
}

func (event *BatchEvent) Merge(other BatchEvent) {
	event.TimeInfo = other.TimeInfo
	event.Trades = append(event.Trades, other.Trades...)
	event.Liquidities = append(event.Liquidities, other.Liquidities...)
}

type EventHandler interface {
	OnEventBatch(event BatchEvent) error
}
