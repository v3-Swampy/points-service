package parsing

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
)

type Client struct {
	interfaces.Provider
}

func NewClient(url string, option ...providers.Option) (*Client, error) {
	provider, err := providers.NewProviderWithOption(url, optionWithDefault(option...))
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to dial %v", url)
	}

	return &Client{
		Provider: provider,
	}, nil
}

func (client *Client) FirstTimestamp() (int64, error) {
	return providers.Call[int64](client.Provider, "firstTimestamp")
}

func (client *Client) LatestTimestamp() (int64, error) {
	return providers.Call[int64](client.Provider, "latestTimestamp")
}

func (client *Client) SnapshotIntervalSecs() (int64, error) {
	return providers.Call[int64](client.Provider, "snapshotInterval")
}

func (client *Client) GetTradeData(pool common.Address, timestamp int64, offset int, limit ...int) (*PagingResult[TradeData], error) {
	if len(limit) == 0 {
		return providers.Call[*PagingResult[TradeData]](client.Provider, "getHourlyTradeData", pool, timestamp, offset)
	}

	return providers.Call[*PagingResult[TradeData]](client.Provider, "getHourlyTradeData", pool, timestamp, offset, limit[0])
}

func (client *Client) GetTradeDataAll(pool common.Address, timestamp int64) ([]TradeData, error) {
	var offset int
	var all []TradeData

	for {
		result, err := client.GetTradeData(pool, timestamp, offset)
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to get trade data with offset %v", offset)
		}

		if result == nil {
			return nil, nil
		}

		offset += len(result.Data)

		if all == nil {
			all = result.Data
		} else {
			all = append(all, result.Data...)
		}

		if len(all) == result.Total {
			return all, nil
		}
	}
}

func (client *Client) GetLiquidityData(pool common.Address, timestamp int64, offset int, limit ...int) (*PagingResult[LiquidityData], error) {
	if len(limit) == 0 {
		return providers.Call[*PagingResult[LiquidityData]](client.Provider, "getHourlyLiquidityData", pool, timestamp, offset)
	}

	return providers.Call[*PagingResult[LiquidityData]](client.Provider, "getHourlyLiquidityData", pool, timestamp, offset, limit[0])
}

func (client *Client) GetLiquidityDataAll(pool common.Address, timestamp int64) ([]LiquidityData, error) {
	var offset int
	var all []LiquidityData

	for {
		result, err := client.GetLiquidityData(pool, timestamp, offset)
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to get liquidity data with offset %v", offset)
		}

		if result == nil {
			return nil, nil
		}

		offset += len(result.Data)

		if all == nil {
			all = result.Data
		} else {
			all = append(all, result.Data...)
		}

		if len(all) == result.Total {
			return all, nil
		}
	}
}
