package sync

import (
	"github.com/shopspring/decimal"
)

type PoolEvent struct {
	Timestamp int64  // unix timestamp in seconds
	User      string // user address, e.g. trader or liquidity provider
	Pool      string // pool address
}

type TradeEvent struct {
	PoolEvent

	Value0 decimal.Decimal // amount0 * price0
	Value1 decimal.Decimal // amount1 * price1
}

type LiquidityEvent struct {
	PoolEvent

	ValueSecs decimal.Decimal // liquidity * price * secs
}

type EventHandler interface {
	OnEventBatch(trades []TradeEvent, liquidities []LiquidityEvent) error
}
