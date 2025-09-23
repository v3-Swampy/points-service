package parsing

import "github.com/ethereum/go-ethereum/common/hexutil"

type PagingResult[T any] struct {
	Total int `json:"total"`
	Data  []T `json:"data"`
}

type TradeData struct {
	UserAddress  string       `json:"user"`
	Token0Volume *hexutil.Big `json:"token0Volume"`
	Token1Volume *hexutil.Big `json:"token1Volume"`
}

type LiquidityData struct {
	UserAddress            string       `json:"user"`
	Token0LiquiditySeconds *hexutil.Big `json:"token0LiquiditySeconds"`
	Token1LiquiditySeconds *hexutil.Big `json:"token1LiquiditySeconds"`
}
