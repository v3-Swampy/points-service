package scan

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Response[T any] struct {
	Status  string
	Message string
	Result  T
}

func (resp *Response[T]) GetResult() (v T, err error) {
	if resp.Status == "1" {
		v = resp.Result
	} else if len(resp.Status) == 0 {
		err = errors.Errorf("Scan Error: %v", resp.Message)
	} else {
		err = errors.Errorf("Scan Error (%v): %v", resp.Status, resp.Message)
	}

	return
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

	return resp.GetResult()
}
