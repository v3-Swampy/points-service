package model

import (
	"strings"

	"github.com/shopspring/decimal"
)

type PagingRequest struct {
	Offset int    `form:"offset" binding:"min=0"`
	Limit  int    `form:"limit" binding:"required,min=1,max=100"`
	Sort   string `form:"sort,default=desc" binding:"omitempty,oneof=asc desc"`
}

func (p *PagingRequest) IsDesc() bool {
	return strings.EqualFold(p.Sort, "desc")
}

type UserPagingRequest struct {
	PagingRequest
	SortField string `form:"sortField,default=trade" binding:"oneof=trade liquidity"`
}

type PoolPagingRequest struct {
	PagingRequest
	SortField string `form:"sortField,default=tvl" binding:"oneof=tvl trade liquidity"`
}

type PagingResult[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

type PagingResultWithUpdatedAt[T any] struct {
	Total     int64  `json:"total"`
	Items     []T    `json:"items"`
	UpdatedAt string `json:"updatedAt"`
}

type UserInfo struct {
	Address         string          `json:"address"`
	TradePoints     decimal.Decimal `json:"tradePoints"`
	LiquidityPoints decimal.Decimal `json:"liquidityPoints"`
}

type PoolInfo struct {
	Address         string          `json:"address"`
	Name            string          `json:"name"`
	Tvl             decimal.Decimal `json:"tvl"`
	TradePoints     decimal.Decimal `json:"tradePoints"`
	LiquidityPoints decimal.Decimal `json:"liquidityPoints"`
}
