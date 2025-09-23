package blockchain

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/blockchain/contract"
)

var ErrSwappiPairNotFound = errors.New("Swappi pair not found in factory")

type PairInfo struct {
	Address common.Address
	Token0  TokenInfo
	Token1  TokenInfo
}

func (info PairInfo) String() string {
	return fmt.Sprintf("%v/%v", info.Token0.Symbol, info.Token1.Symbol)
}

type SwappiAddresses struct {
	Factory common.Address
	USDT    common.Address
	WCFX    common.Address
}

type Swappi struct {
	caller    bind.ContractCaller
	erc20     *ERC20
	addresses SwappiAddresses
	cache     sync.Map
}

func NewSwappi(caller bind.ContractCaller, erc20 *ERC20, addresses SwappiAddresses) *Swappi {
	return &Swappi{
		caller:    caller,
		erc20:     erc20,
		addresses: addresses,
	}
}

// GetPairInfo retrieves pair token info from blockchain or returns the cached value.
func (swappi *Swappi) GetPairInfo(pair common.Address) (PairInfo, error) {
	// check cache at first
	if val, ok := swappi.cache.Load(pair); ok {
		return val.(PairInfo), nil
	}

	info, err := swappi.getPairInfo(pair)
	if err != nil {
		return PairInfo{}, err
	}

	// cache value
	swappi.cache.Store(pair, info)

	return info, nil
}

func (swappi *Swappi) getPairInfo(pair common.Address) (PairInfo, error) {
	info := PairInfo{
		Address: pair,
	}

	// retrieves token0 & token1 from pair token
	pairCaller, err := contract.NewSwappiPairCaller(pair, swappi.caller)
	if err != nil {
		return PairInfo{}, errors.WithMessage(err, "Failed to create Pair caller")
	}

	token0, err := pairCaller.Token0(nil)
	if err != nil {
		return PairInfo{}, errors.WithMessage(err, "Failed to query token0 of LP token")
	}

	token1, err := pairCaller.Token1(nil)
	if err != nil {
		return PairInfo{}, errors.WithMessage(err, "Failed to query token1 of LP token")
	}

	// retrieves token info of token0 and token1
	if info.Token0, err = swappi.erc20.GetTokenInfo(token0); err != nil {
		return PairInfo{}, errors.WithMessage(err, "Failed to query token0 info")
	}

	if info.Token1, err = swappi.erc20.GetTokenInfo(token1); err != nil {
		return PairInfo{}, errors.WithMessage(err, "Failed to query token1 info")
	}

	return info, nil
}

// GetTokenPrice calculates the price of given baseToken and quoteToken.
//
// It returns ErrSwappiPairNotFound if pair not found in given factory.
func (swappi *Swappi) GetTokenPrice(opts *bind.CallOpts, baseToken, quoteToken common.Address) (decimal.Decimal, error) {
	// get pair from factory
	factoryCaller, err := contract.NewSwappiFactoryCaller(swappi.addresses.Factory, swappi.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Factory caller")
	}

	pair, err := factoryCaller.GetPair(opts, baseToken, quoteToken)
	if err != nil {
		return decimal.Zero, errors.WithMessagef(err, "Failed to get pair from factory")
	}

	if pair.Cmp(common.Address{}) == 0 {
		return decimal.Zero, ErrSwappiPairNotFound
	}

	info, err := swappi.GetPairInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair info")
	}

	// get reserves from pair
	pairCaller, err := contract.NewSwappiPairCaller(pair, swappi.caller)
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

// GetTokenPriceRouted calculates the price of given routes.
//
// It returns ErrSwappiPairNotFound if any pair not found in the routes or the number of given routes is less than 2.
func (swappi *Swappi) GetTokenPriceRouted(opts *bind.CallOpts, routes ...common.Address) (decimal.Decimal, error) {
	numTokens := len(routes)
	if numTokens < 2 {
		return decimal.Zero, ErrSwappiPairNotFound
	}

	price := decimal.New(1, 0)

	for i := 1; i < numTokens; i++ {
		midPrice, err := swappi.GetTokenPrice(opts, routes[i-1], routes[i])
		if err != nil {
			return decimal.Zero, errors.WithMessagef(err, "Failed to get mid price, i = %v", i)
		}

		price = price.Mul(midPrice)
	}

	return price, nil
}

// GetTokenPriceAuto calculates the price of give token via token/WCFX/USDT or token/USDT.
//
// It returns ErrSwappiPairNotFound if relevant pair not found.
func (swappi *Swappi) GetTokenPriceAuto(opts *bind.CallOpts, token common.Address) (decimal.Decimal, error) {
	if token == swappi.addresses.USDT {
		return decimal.NewFromInt(1), nil
	}

	if token == swappi.addresses.WCFX {
		return swappi.GetTokenPrice(opts, swappi.addresses.WCFX, swappi.addresses.USDT)
	}

	// try to get price by token/WCFX/USDT with priority
	price, err := swappi.GetTokenPriceRouted(opts, token, swappi.addresses.WCFX, swappi.addresses.USDT)
	if err == nil {
		return price, nil
	}

	if err != ErrSwappiPairNotFound {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/WCFX/USDT")
	}

	// otherwise, try to get price by token/USDT
	price, err = swappi.GetTokenPrice(opts, token, swappi.addresses.USDT)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/USDT")
	}

	return price, nil
}

// GetPairTVL calculates the TVL of given pair via reserves.
func (swappi *Swappi) GetPairTVL(opts *bind.CallOpts, pair common.Address) (decimal.Decimal, error) {
	info, err := swappi.GetPairInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair info")
	}

	// get prices
	price0, err := swappi.GetTokenPriceAuto(opts, info.Token0.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token0")
	}

	price1, err := swappi.GetTokenPriceAuto(opts, info.Token1.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token1")
	}

	// get reserves from pair
	pairCaller, err := contract.NewSwappiPairCaller(pair, swappi.caller)
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

// GetPairTokenPrice calculates the price of give pair token.
func (swappi *Swappi) GetPairTokenPrice(opts *bind.CallOpts, pair common.Address) (decimal.Decimal, error) {
	info, err := swappi.erc20.GetTokenInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair token info")
	}

	pairCaller, err := contract.NewSwappiPairCaller(pair, swappi.caller)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to create Pair caller")
	}

	supply, err := pairCaller.TotalSupply(opts)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get total supply of pair")
	}

	supplyDecimal := decimal.NewFromBigInt(supply, -int32(info.Decimals))

	tvl, err := swappi.GetPairTVL(opts, pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair TVL")
	}

	return tvl.Div(supplyDecimal), nil
}

// GetPairTVL calculates the TVL of given pair via tokens balances hold by pool.
func (swappi *Swappi) GetPairTVLByBalances(opts *bind.CallOpts, pair common.Address) (decimal.Decimal, error) {
	info, err := swappi.GetPairInfo(pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pair info")
	}

	// get prices
	price0, err := swappi.GetTokenPriceAuto(opts, info.Token0.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token0")
	}

	price1, err := swappi.GetTokenPriceAuto(opts, info.Token1.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token1")
	}

	// get balances of pair
	balance0, err := swappi.erc20.GetBalance(opts, info.Token0.Address, pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token0 balance of pair")
	}

	balance1, err := swappi.erc20.GetBalance(opts, info.Token1.Address, pair)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get token1 balance of pair")
	}

	value0 := balance0.Mul(price0)
	value1 := balance1.Mul(price1)

	return value0.Add(value1), nil
}
