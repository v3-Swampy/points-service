package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var ErrVswapPoolNotFound = errors.New("vSwap pool not found to calculate price")

type Vswap struct {
	swappi *Swappi

	wcfxUsdtPool common.Address
}

func NewVswap(swappi *Swappi, wcfxUsdtPool common.Address) *Vswap {
	return &Vswap{swappi, wcfxUsdtPool}
}

func (vswap *Vswap) GetPoolInfo(pool common.Address) (PairInfo, error) {
	return vswap.swappi.GetPairInfo(pool)
}

// GetTokenPrice calculates the price of given token in pool.
//
// It will returns error if the given token not found in pool.
//
// Note, it returns 0 if pool balance of any token is 0.
func (vswap *Vswap) GetTokenPrice(opts *bind.CallOpts, pool, token common.Address) (decimal.Decimal, error) {
	// get pool info
	info, err := vswap.GetPoolInfo(pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool info")
	}

	// get balances
	balance0, err := vswap.swappi.erc20.GetBalance(opts, info.Token0.Address, pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool balance of token0")
	}

	balance1, err := vswap.swappi.erc20.GetBalance(opts, info.Token1.Address, pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool balance of token1")
	}

	if balance0.IsZero() || balance1.IsZero() {
		return decimal.Zero, nil
	}

	// calculate price
	if token == info.Token0.Address {
		return balance1.Div(balance0), nil
	}

	if token == info.Token1.Address {
		return balance0.Div(balance1), nil
	}

	return decimal.Zero, errors.Errorf("Token not found in pool %v", info)
}

// GetTokenPrice calculates the USDT price of given token in pool.
//
// It will returns error if the given token not found in pool.
//
// Note, it returns 0 if pool balance of any token is 0.
func (vswap *Vswap) GetTokenPriceUSDT(opts *bind.CallOpts, pool, token common.Address) (decimal.Decimal, error) {
	if token == vswap.swappi.addresses.USDT {
		return decimal.NewFromInt(1), nil
	}

	if token == vswap.swappi.addresses.WCFX {
		return vswap.GetTokenPrice(opts, vswap.wcfxUsdtPool, token)
	}

	// get pool info
	info, err := vswap.GetPoolInfo(pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool info")
	}

	var other common.Address
	switch token {
	case info.Token0.Address:
		other = info.Token1.Address
	case info.Token1.Address:
		other = info.Token0.Address
	default:
		return decimal.Zero, errors.Errorf("Token not found in pool %v", info)
	}

	// token/usdt
	if other == vswap.swappi.addresses.USDT {
		return vswap.GetTokenPrice(opts, pool, token)
	}

	if other != vswap.swappi.addresses.WCFX {
		return decimal.Zero, ErrVswapPoolNotFound
	}

	// token/wcfx/usdt
	wcfxPrice, err := vswap.GetTokenPrice(opts, pool, token)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to calculate token price by token/WCFX")
	}

	if wcfxPrice.IsZero() {
		return decimal.Zero, nil
	}

	usdtPrice, err := vswap.GetTokenPrice(opts, vswap.wcfxUsdtPool, vswap.swappi.addresses.WCFX)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to calcuate WCFX price by WCFX/USDT")
	}

	return wcfxPrice.Mul(usdtPrice), nil
}

func (vswap *Vswap) GetPoolTVL(opts *bind.CallOpts, pool common.Address) (decimal.Decimal, error) {
	info, err := vswap.GetPoolInfo(pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool info")
	}

	// get balances
	balance0, err := vswap.swappi.erc20.GetBalance(opts, info.Token0.Address, pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool balance of token0")
	}

	balance1, err := vswap.swappi.erc20.GetBalance(opts, info.Token1.Address, pool)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get pool balance of token1")
	}

	if balance0.IsZero() || balance1.IsZero() {
		return decimal.Zero, nil
	}

	// get prices
	price0, err := vswap.GetTokenPriceUSDT(opts, pool, info.Token0.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token0")
	}

	price1, err := vswap.GetTokenPriceUSDT(opts, pool, info.Token1.Address)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price of token1")
	}

	// calculate TVL
	tvl0 := balance0.Mul(price0)
	tvl1 := balance1.Mul(price1)

	return tvl0.Add(tvl1), nil
}
