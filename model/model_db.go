package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/sync/blockchain"
)

var Tables = []any{&User{}, &Pool{}}

type Model struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
}

type User struct {
	Model
	Address         string          `gorm:"size:64;not null;unique" json:"address"`
	TradePoints     decimal.Decimal `gorm:"type:decimal(10,0);not null;index:idx_trade_points" json:"tradePoints"`
	LiquidityPoints decimal.Decimal `gorm:"type:decimal(11,1);not null;index:idx_liquidity_points" json:"liquidityPoints"`
}

func NewUser(address string, tradePoints decimal.Decimal, liquidityPoints decimal.Decimal, time time.Time) *User {
	return &User{
		Address:         address,
		TradePoints:     tradePoints,
		LiquidityPoints: liquidityPoints,
		Model: Model{
			CreatedAt: time,
			UpdatedAt: time,
		},
	}
}

type Pool struct {
	Model
	Address         string          `gorm:"size:64;not null;unique" json:"address"`
	Token0          string          `gorm:"size:64;not null" json:"token0"`
	Token1          string          `gorm:"size:64;not null" json:"token1"`
	Tvl             decimal.Decimal `gorm:"type:decimal(10,0);not null;index:idx_tvl" json:"tvl"`
	TradePoints     decimal.Decimal `gorm:"type:decimal(10,0);not null;index:idx_trade_points" json:"tradePoints"`
	LiquidityPoints decimal.Decimal `gorm:"type:decimal(11,1);not null;index:idx_liquidity_points" json:"liquidityPoints"`

	Token0Name     string `gorm:"size:128" json:"token0Name"`
	Token0Symbol   string `gorm:"size:128" json:"token0Symbol"`
	Token0Decimals uint8  `gorm:"" json:"token0Decimals"`
	Token1Name     string `gorm:"size:128" json:"token1Name"`
	Token1Symbol   string `gorm:"size:128" json:"token1Symbol"`
	Token1Decimals uint8  `gorm:"" json:"token1Decimals"`
}

func NewPool(pool blockchain.PairInfo, tradePoints decimal.Decimal, liquidityPoints decimal.Decimal, time time.Time) *Pool {
	return &Pool{
		Address:         pool.Address.String(),
		Token0:          pool.Token0.Address.String(),
		Token1:          pool.Token1.Address.String(),
		Tvl:             decimal.Zero, //TODO
		TradePoints:     tradePoints,
		LiquidityPoints: liquidityPoints,

		Token0Name:     pool.Token0.Name,
		Token0Symbol:   pool.Token0.Symbol,
		Token0Decimals: pool.Token0.Decimals,
		Token1Name:     pool.Token1.Name,
		Token1Symbol:   pool.Token1.Symbol,
		Token1Decimals: pool.Token1.Decimals,

		Model: Model{
			CreatedAt: time,
			UpdatedAt: time,
		},
	}
}
