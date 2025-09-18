package parsing

import "github.com/ethereum/go-ethereum/common/hexutil"

type PagingResult[T any] struct {
	Total int `json:"total"`
	Data  []T `json:"data"`
}

type TradeData struct {
	UserAddress   string       `json:"userAddress"`
	Token0Volumes *hexutil.Big `json:"token0volumes"`
	Token1Volumes *hexutil.Big `json:"token1volumes"`
}

type LiquidityData struct {
	UserAddress      string       `json:"userAddress"`
	LiquiditySeconds *hexutil.Big `json:"liquiditySeconds"` // liquidity * seconds
}
