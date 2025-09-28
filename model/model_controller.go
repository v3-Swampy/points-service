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
	Total     int64 `json:"total"`
	Items     []T   `json:"items"`
	UpdatedAt int64 `json:"updatedAt"`
}

type UserInfo struct {
	Address         string          `json:"address"`
	TradePoints     decimal.Decimal `json:"tradePoints"`
	LiquidityPoints decimal.Decimal `json:"liquidityPoints"`
}

type PoolParamInfo struct {
	Address         string          `json:"address"`
	Token0          string          `json:"token0"`
	Token1          string          `json:"token1"`
	Token0Symbol    string          `json:"token0Symbol"`
	Token1Symbol    string          `json:"token1Symbol"`
	TradeWeight     decimal.Decimal `json:"tradeWeight"`
	LiquidityWeight decimal.Decimal `json:"liquidityWeight"`
}

type PoolInfo struct {
	PoolParamInfo
	Fee uint32          `json:"fee"`
	Tvl decimal.Decimal `json:"tvl"`
}
