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

func (bc *Blockchain) GetSwappiTokenPriceAuto(opts *bind.CallOpts, token common.Address, addresses SwappiAddresses) (decimal.Decimal, error) {
	if token == addresses.USDT {
		return decimal.NewFromInt(1), nil
	}

	if token == addresses.WCFX {
		return bc.GetSwappiTokenPrice(opts, addresses.Factory, addresses.WCFX, addresses.USDT)
	}

	// try to get price by token/USDT
	price, err := bc.GetSwappiTokenPrice(opts, addresses.Factory, token, addresses.USDT)
	if err == nil {
		return price, nil
	}

	if err != ErrSwappiPairNotFound {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/USDT")
	}

	// otherwise, try to get price by token/WCFX/USDT
	wcfxPrice, err := bc.GetSwappiTokenPrice(opts, addresses.Factory, token, addresses.WCFX)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by token/WCFX")
	}

	usdtPrice, err := bc.GetSwappiTokenPrice(opts, addresses.Factory, addresses.WCFX, addresses.USDT)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get price by WCFX/USDT")
	}

	return wcfxPrice.Mul(usdtPrice), nil
}

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

func (bc *Blockchain) GetSwappiTokenPriceLP(opts *bind.CallOpts, pair common.Address, addresses SwappiAddresses) (decimal.Decimal, error) {
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

	// calculate LP token price
	supply, err := pairCaller.TotalSupply(opts)
	if err != nil {
		return decimal.Zero, errors.WithMessage(err, "Failed to get total supply of pair")
	}

	supplyDecimal := decimal.NewFromBigInt(supply, -int32(info.TokenLP.Decimals))

	return value0.Add(value1).Div(supplyDecimal), nil
}
