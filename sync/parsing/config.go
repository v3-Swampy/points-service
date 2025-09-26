package parsing

import "github.com/mcuadros/go-defaults"

type Config struct {
	Poller struct {
		RpcUrl  string
		ScanUrl string
		Option  PollOption
	}

	Emitter EmitOption
	Batcher BatchOption
}

func optionWithDefault[T any](option ...T) T {
	var opt T

	if len(option) > 0 {
		opt = option[0]
	}

	defaults.SetDefaults(&opt)

	return opt
}
