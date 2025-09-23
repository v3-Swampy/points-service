package blockchain

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	URL    string
	Scan   string
	Swappi SwappiConfig
	Vswap  VswapConfig
}

type SwappiConfig struct {
	Factory string
	USDT    string
	WCFX    string
}

type VswapConfig struct {
	WcfxUsdtPool string
}

func (config *SwappiConfig) ToAddresses() SwappiAddresses {
	return SwappiAddresses{
		Factory: common.HexToAddress(config.Factory),
		USDT:    common.HexToAddress(config.USDT),
		WCFX:    common.HexToAddress(config.WCFX),
	}
}
