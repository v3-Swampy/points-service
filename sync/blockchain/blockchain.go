package blockchain

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/sync/blockchain/contract"
)

type TokenInfo struct {
	Address  common.Address
	Name     string
	Symbol   string
	Decimals uint8
}

type PoolInfo struct {
	TokenLP TokenInfo
	Token0  TokenInfo
	Token1  TokenInfo
}

type Blockchain struct {
	caller         bind.ContractCaller
	tokenInfoCache sync.Map
	poolInfoCache  sync.Map
}

func NewBlockchain(client *web3go.Client) *Blockchain {
	caller, _ := client.ToClientForContract()

	return &Blockchain{
		caller: caller,
	}
}

func (bc *Blockchain) GetTokenInfo(token common.Address) (TokenInfo, error) {
	if val, ok := bc.tokenInfoCache.Load(token); ok {
		return val.(TokenInfo), nil
	}

	caller, err := contract.NewERC20Caller(token, bc.caller)
	if err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to create ERC20 caller")
	}

	info := TokenInfo{
		Address: token,
	}

	if info.Name, err = caller.Name(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token name")
	}

	if info.Symbol, err = caller.Symbol(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token symbol")
	}

	if info.Decimals, err = caller.Decimals(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token decimals")
	}

	bc.tokenInfoCache.Store(token, info)

	return info, nil
}

func (bc *Blockchain) GetPairTokenInfo(pairToken common.Address) (PoolInfo, error) {
	if val, ok := bc.poolInfoCache.Load(pairToken); ok {
		return val.(PoolInfo), nil
	}

	var info PoolInfo

	pairCaller, err := contract.NewSwappiPairCaller(pairToken, bc.caller)
	if err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to create Pair caller")
	}

	token0, err := pairCaller.Token0(nil)
	if err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token0 of LP token")
	}

	token1, err := pairCaller.Token1(nil)
	if err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token1 of LP token")
	}

	if info.TokenLP, err = bc.GetTokenInfo(pairToken); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query LP token info")
	}

	if info.Token0, err = bc.GetTokenInfo(token0); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token0 info")
	}

	if info.Token1, err = bc.GetTokenInfo(token1); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token1 info")
	}

	bc.poolInfoCache.Store(pairToken, info)

	return info, nil
}

func (bc *Blockchain) GetTokenPrice(token common.Address) (decimal.Decimal, error) {
	// TODO
	return decimal.NewFromFloat(1), nil
}

func (bc *Blockchain) GetTokenPriceLP(lpToken common.Address) (decimal.Decimal, error) {
	// TODO
	return decimal.NewFromFloat(2), nil
}
