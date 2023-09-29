package prober

import (
	"time"

	"crdx.org/lighthouse/service"
)

// Documentation can be found attached to the definition of Config in service/service.go.

type config struct{}

func (config) Name() string {
	return "prober"
}

func (config) InitialRestartInterval() time.Duration {
	return time.Second
}

func (config) NextRestartInterval(previous time.Duration) time.Duration {
	return service.ExponentialBackoff(previous)
}

func (config) RunInterval() time.Duration {
	return time.Hour
}

func (config) InitialStartDelay() time.Duration {
	return time.Hour
}
