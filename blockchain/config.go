package blockchain

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	URL    string
	Swappi SwappiConfig
}

type SwappiConfig struct {
	Factory string
	USDT    string
	WCFX    string
}

func (config *SwappiConfig) ToAddresses() SwappiAddresses {
	return SwappiAddresses{
		Factory: common.HexToAddress(config.Factory),
		USDT:    common.HexToAddress(config.USDT),
		WCFX:    common.HexToAddress(config.WCFX),
	}
}
