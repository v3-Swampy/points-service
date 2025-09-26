package scan

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/pkg/errors"
)

type Option struct {
	DebugEnabled   bool
	RequestTimeout time.Duration `default:"3s"`
}

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
}

func NewApi(url string, option ...Option) *Api {
	var opt Option
	if len(option) > 0 {
		opt = option[0]
	}

	defaults.SetDefaults(&opt)

	return &Api{
		client: resty.New().
			SetBaseURL(url).
			SetDebug(opt.DebugEnabled).
			SetTimeout(opt.RequestTimeout),
	}
}

func (api *Api) GetBlockNumberByTime(timestampSecs int64, after bool) (uint64, error) {
	closest := "before"
	if after {
		closest = "after"
	}

	url := fmt.Sprintf("/api?module=block&action=getblocknobytime&timestamp=%v&closest=%v", timestampSecs, closest)

	var resp Response[uint64]

	if _, err := api.client.R().SetResult(&resp).Get(url); err != nil {
		return 0, err
	}

	return resp.GetResult()
}
