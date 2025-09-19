package blockchain

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/v3-Swampy/points-service/sync/blockchain/contract"
)

type TokenInfo struct {
	Address  common.Address
	Name     string
	Symbol   string
	Decimals uint8
}

type ERC20 struct {
	caller bind.ContractCaller
	cache  sync.Map
}

func NewERC20(caller bind.ContractCaller) *ERC20 {
	return &ERC20{
		caller: caller,
	}
}

// GetTokenInfo retrieves ERC20 token info from blockchain or returns the cached value.
func (erc20 *ERC20) GetTokenInfo(token common.Address) (TokenInfo, error) {
	// check cache at first
	if val, ok := erc20.cache.Load(token); ok {
		return val.(TokenInfo), nil
	}

	caller, err := contract.NewERC20Caller(token, erc20.caller)
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
	erc20.cache.Store(token, info)

	return info, nil
}
