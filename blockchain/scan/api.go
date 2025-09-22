package scan

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Response[T any] struct {
	Status  string
	Message string
	Result  T
}

type Api struct {
	client *resty.Client
	url    string
}

func NewApi(url string) *Api {
	return &Api{
		client: resty.New(),
		url:    url,
	}
}

func (api *Api) GetBlockNumberByTime(timestampSecs int64, after bool) (uint64, error) {
	closest := "before"
	if after {
		closest = "after"
	}

	url := fmt.Sprintf("%v/api?module=block&action=getblocknobytime&timestamp=%v&closest=%v", api.url, timestampSecs, closest)

	var resp Response[uint64]

	if _, err := api.client.R().SetResult(&resp).Get(url); err != nil {
		return 0, err
	}

	return resp.Result, nil
}
