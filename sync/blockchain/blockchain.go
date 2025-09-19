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

var ErrSwappiPairNotFound = errors.New("Pair not found in factory")

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

type SwappiAddresses struct {
	Factory common.Address
	USDT    common.Address
	WCFX    common.Address
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

// GetTokenInfo retrieves ERC20 token info from blockchain or returns the cached value.
func (bc *Blockchain) GetTokenInfo(token common.Address) (TokenInfo, error) {
	// check cache at first
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

	// retrieves name, symbol and decimals
	if info.Name, err = caller.Name(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token name")
	}

	if info.Symbol, err = caller.Symbol(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token symbol")
	}

	if info.Decimals, err = caller.Decimals(nil); err != nil {
		return TokenInfo{}, errors.WithMessage(err, "Failed to query token decimals")
	}

	// cache value
	bc.tokenInfoCache.Store(token, info)

	return info, nil
}

// GetPairTokenInfo retrieves LP token info from blockchain or returns the cached value.
func (bc *Blockchain) GetPairTokenInfo(pairToken common.Address) (PoolInfo, error) {
	// check cache at first
	if val, ok := bc.poolInfoCache.Load(pairToken); ok {
		return val.(PoolInfo), nil
	}

	var info PoolInfo

	// retrieves token0 & token1 from pair token
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

	// retrieves token info of LP, token0 and token1
	if info.TokenLP, err = bc.GetTokenInfo(pairToken); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query LP token info")
	}

	if info.Token0, err = bc.GetTokenInfo(token0); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token0 info")
	}

	if info.Token1, err = bc.GetTokenInfo(token1); err != nil {
		return PoolInfo{}, errors.WithMessage(err, "Failed to query token1 info")
	}

	// cache value
	bc.poolInfoCache.Store(pairToken, info)

	return info, nil
}

// GetSwappiTokenPrice calculates the Swappi price of given baseToken and quoteToken.
//
// It returns ErrSwappiPairNotFound if pair not found in given Swappi factory.
func (bc *Blockchain) GetSwappiTokenPrice(opts *bind.CallOpts, swappiFactory, baseToken, quoteToken common.Address) (decimal.Decimal, error) {
	// get pair token from factory
	factoryCaller, err := contract.NewSwappiFactoryCaller(swappiFactory, bc.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Factory caller")
	}

	pair, err := factoryCaller.GetPair(opts, baseToken, quoteToken)
	if err != nil {
		return decimal.Zero, errors.WithMessagef(err, "Failed to get pair from Swappi factory")
	}

	if pair.Cmp(common.Address{}) == 0 {
		return decimal.Zero, ErrSwappiPairNotFound
	}

	info, err := bc.GetPairTokenInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair token info")
	}

	// get reserves from pair
	pairCaller, err := contract.NewSwappiPairCaller(pair, bc.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Pair caller")
	}

	reserves, err := pairCaller.GetReserves(opts)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get reserves from pair")
	}

	reserve0 := decimal.NewFromBigInt(reserves.Reserve0, -int32(info.Token0.Decimals))
	reserve1 := decimal.NewFromBigInt(reserves.Reserve1, -int32(info.Token1.Decimals))

	if info.Token0.Address == baseToken {
		return reserve1.Div(reserve0), nil
	}

	return reserve0.Div(reserve1), nil
}

// GetSwappiTokenPriceRouted calculates the Swappi price of given routes.
//
// It returns ErrSwappiPairNotFound if any pair not found in the routes or the number of given routes is less than 2.
func (bc *Blockchain) GetSwappiTokenPriceRouted(opts *bind.CallOpts, swappiFactory common.Address, routes ...common.Address) (decimal.Decimal, error) {
	numTokens := len(routes)
	if numTokens < 2 {
		return decimal.Zero, ErrSwappiPairNotFound
	}

	result := decimal.New(1, 0)

	for i := 1; i < numTokens; i++ {
		price, err := bc.GetSwappiTokenPrice(opts, swappiFactory, routes[i-1], routes[i])
		if err != nil {
			return decimal.Zero, errors.WithMessagef(err, "Failed to get price, i = %v", i)
		}

		result = result.Mul(price)
	}

	return result, nil
}

// GetSwappiTokenPriceAuto calculates the Swappi price of give token via token/CFX/USDT or token/USDT.
//
// It returns ErrSwappiPairNotFound if relevant pair not found.
func (bc *Blockchain) GetSwappiTokenPriceAuto(opts *bind.CallOpts, token common.Address, addresses SwappiAddresses) (decimal.Decimal, error) {
	if token == addresses.USDT {
		return decimal.NewFromInt(1), nil
	}

	if token == addresses.WCFX {
		return bc.GetSwappiTokenPrice(opts, addresses.Factory, addresses.WCFX, addresses.USDT)
	}

	// try to get price by token/WCFX/USDT with priority
	price, err := bc.GetSwappiTokenPriceRouted(opts, addresses.Factory, token, addresses.WCFX, addresses.USDT)
	if err == nil {
		return price, nil
	}

	if err != ErrSwappiPairNotFound {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/WCFX/USDT")
	}

	// otherwise, try to get price by token/USDT
	price, err = bc.GetSwappiTokenPrice(opts, addresses.Factory, token, addresses.USDT)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/USDT")
	}

	return price, nil
}

// GetSwappiTokenPriceLP calculates the Swappi price of give LP token.
func (bc *Blockchain) GetSwappiTokenPriceLP(opts *bind.CallOpts, pair common.Address, addresses SwappiAddresses) (decimal.Decimal, error) {
	info, err := bc.GetPairTokenInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair token info")
	}

	pairCaller, err := contract.NewSwappiPairCaller(pair, bc.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Pair caller")
	}

	supply, err := pairCaller.TotalSupply(opts)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get total supply of pair")
	}

	supplyDecimal := decimal.NewFromBigInt(supply, -int32(info.TokenLP.Decimals))

	tvl, err := bc.GetSwappiPoolTVL(opts, pair, addresses)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get TVL of pool")
	}

	return tvl.Div(supplyDecimal), nil
}

// GetSwappiPoolTVL calculates the Swappi TVL of given pool.
func (bc *Blockchain) GetSwappiPoolTVL(opts *bind.CallOpts, pair common.Address, addresses SwappiAddresses) (decimal.Decimal, error) {
	info, err := bc.GetPairTokenInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair token info")
	}

	// get prices
	price0, err := bc.GetSwappiTokenPriceAuto(opts, info.Token0.Address, addresses)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token0")
	}

	price1, err := bc.GetSwappiTokenPriceAuto(opts, info.Token1.Address, addresses)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token1")
	}

	// get reserves from pair
	pairCaller, err := contract.NewSwappiPairCaller(pair, bc.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Pair caller")
	}

	reserves, err := pairCaller.GetReserves(opts)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get reserves from pair")
	}

	value0 := decimal.NewFromBigInt(reserves.Reserve0, -int32(info.Token0.Decimals)).Mul(price0)
	value1 := decimal.NewFromBigInt(reserves.Reserve1, -int32(info.Token1.Decimals)).Mul(price1)

	return value0.Add(value1), nil
}
