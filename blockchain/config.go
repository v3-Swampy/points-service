package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

type Config struct {
	URL    string
	Option providers.Option

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
