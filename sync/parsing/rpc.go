package parsing

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
)

type Client struct {
	interfaces.Provider
}

func NewClient(url string) (*Client, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to dial %v", url)
	}

	return &Client{
		Provider: client,
	}, nil
}

func (client *Client) GetHourlyTradeData(ctx context.Context, pool string, hourTimestamp int64, offset int, limit ...int) (*PagingResult[TradeData], error) {
	if len(limit) == 0 {
		return providers.CallContext[*PagingResult[TradeData]](client.Provider, ctx, "parser_getHourlyTradeData", pool, hourTimestamp, offset)
	}

	return providers.CallContext[*PagingResult[TradeData]](client.Provider, ctx, "parser_getHourlyTradeData", pool, hourTimestamp, offset, limit[0])
}

func (client *Client) GetHourlyTradeDataAll(ctx context.Context, pool string, hourTimestamp int64) ([]TradeData, error) {
	var offset int
	var all []TradeData

	for {
		result, err := client.GetHourlyTradeData(ctx, pool, hourTimestamp, offset)
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to get hourly trade data with offset %v", offset)
		}

		if result == nil {
			return nil, nil
		}

		offset += len(result.Data)
		all = append(all, result.Data...)

		if len(all) == result.Total {
			return all, nil
		}
	}
}

func (client *Client) GetHourlyLiquidityData(ctx context.Context, pool string, hourTimestamp int64, offset int, limit ...int) (*PagingResult[LiquidityData], error) {
	if len(limit) == 0 {
		return providers.CallContext[*PagingResult[LiquidityData]](client.Provider, ctx, "parser_getHourlyLiquidityData", pool, hourTimestamp, offset)
	}

	return providers.CallContext[*PagingResult[LiquidityData]](client.Provider, ctx, "parser_getHourlyLiquidityData", pool, hourTimestamp, offset, limit[0])
}

func (client *Client) GetHourlyLiquidityDataAll(ctx context.Context, pool string, hourTimestamp int64) ([]LiquidityData, error) {
	var offset int
	var all []LiquidityData

	for {
		result, err := client.GetHourlyLiquidityData(ctx, pool, hourTimestamp, offset)
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to get hourly liquidity data with offset %v", offset)
		}

		if result == nil {
			return nil, nil
		}

		offset += len(result.Data)
		all = append(all, result.Data...)

		if len(all) == result.Total {
			return all, nil
		}
	}
}
