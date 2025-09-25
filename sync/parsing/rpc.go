package parsing

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
)

type Client struct {
	interfaces.Provider
}

func NewClient(url string) (*Client, error) {
	option := providers.Option{
		RequestTimeout: 5 * time.Second,
	}

	provider, err := providers.NewProviderWithOption(url, option)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to dial %v", url)
	}

	return &Client{
		Provider: provider,
	}, nil
}

func (client *Client) FirstTimestamp(ctx context.Context) (int64, error) {
	return providers.CallContext[int64](client.Provider, ctx, "firstTimestamp")
}

func (client *Client) LatestTimestamp(ctx context.Context) (int64, error) {
	return providers.CallContext[int64](client.Provider, ctx, "latestTimestamp")
}

func (client *Client) SnapshotIntervalSecs(ctx context.Context) (int64, error) {
	return providers.CallContext[int64](client.Provider, ctx, "snapshotInterval")
}

func (client *Client) GetTradeData(ctx context.Context, pool common.Address, timestamp int64, offset int, limit ...int) (*PagingResult[TradeData], error) {
	if len(limit) == 0 {
		return providers.CallContext[*PagingResult[TradeData]](client.Provider, ctx, "getHourlyTradeData", pool, timestamp, offset)
	}

	return providers.CallContext[*PagingResult[TradeData]](client.Provider, ctx, "getHourlyTradeData", pool, timestamp, offset, limit[0])
}

func (client *Client) GetTradeDataAll(ctx context.Context, pool common.Address, timestamp int64) ([]TradeData, error) {
	var offset int
	var all []TradeData

	for {
		result, err := client.GetTradeData(ctx, pool, timestamp, offset)
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

func (client *Client) GetLiquidityData(ctx context.Context, pool common.Address, timestamp int64, offset int, limit ...int) (*PagingResult[LiquidityData], error) {
	if len(limit) == 0 {
		return providers.CallContext[*PagingResult[LiquidityData]](client.Provider, ctx, "getHourlyLiquidityData", pool, timestamp, offset)
	}

	return providers.CallContext[*PagingResult[LiquidityData]](client.Provider, ctx, "getHourlyLiquidityData", pool, timestamp, offset, limit[0])
}

func (client *Client) GetLiquidityDataAll(ctx context.Context, pool common.Address, timestamp int64) ([]LiquidityData, error) {
	var offset int
	var all []LiquidityData

	for {
		result, err := client.GetLiquidityData(ctx, pool, timestamp, offset)
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
