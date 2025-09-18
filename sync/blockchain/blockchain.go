package blockchain

import (
	"github.com/openweb3/web3go"
	"github.com/shopspring/decimal"
)

type TokenInfo struct {
	Address  string
	Decimals uint8
}

type PoolInfo struct {
	TokenLP TokenInfo
	Token0  TokenInfo
	Token1  TokenInfo
}

type Blockchain struct {
	client *web3go.Client
}

func NewBlockchain(client *web3go.Client) *Blockchain {
	return &Blockchain{client}
}

func (bc *Blockchain) GetPoolInfo(poolAddress string) (PoolInfo, error) {
	// TODO cacheable
	return PoolInfo{}, nil
}

func (bc *Blockchain) GetTokenPrice(tokenAddress string) (decimal.Decimal, error) {
	// TODO
	return decimal.NewFromFloat(1), nil
}

func (bc *Blockchain) GetTokenPriceLP(lpTokenAddress string) (decimal.Decimal, error) {
	// TODO
	return decimal.NewFromFloat(2), nil
}
