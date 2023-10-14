package services

import (
	"time"

	"log/slog"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
	"crdx.org/lighthouse/util"
)

type Args struct {
	Logger *slog.Logger
}

type Service interface {
	// Init initialises the service and returns an error if init failed. If an error is returned
	// then the service will not run.
	Init(*Args) error

	// Run runs the service and returns an error for an unrecoverable failure.
	Run() error
}

type Config struct {
	// Service is an instance of the service.
	Service Service

	// StartDelay is the delay before starting the service.
	StartDelay time.Duration

	// RunInterval is the interval between runs.
	RunInterval time.Duration

	// InitialRestartInterval is the initial restart interval.
	InitialRestartInterval time.Duration

	// NextRestartInterval is a function that receives the last restart interval and returns the
	// next one.
	NextRestartInterval func(time.Duration) time.Duration
}

// applyDefaults updates config by setting default values for any unset fields, ensuring that the
// config object is valid.
func applyDefaults(config *Config) {
	if config.InitialRestartInterval == 0 {
		config.InitialRestartInterval = time.Second
	}

	if config.NextRestartInterval == nil {
		config.NextRestartInterval = ExponentialBackoff
	}
}

func Start(name string, config *Config) {
	applyDefaults(config)

	log := logger.With("service", name)

	if err := config.Service.Init(&Args{Logger: log}); err != nil {
		log.Error("service init failed", "msg", err)
		return
	}

	log.Info(
		"service init complete",
		"run_interval", config.RunInterval.String(),
		"initial_delay", config.StartDelay.String(),
	)

	go func() {
		if config.StartDelay > 0 {
			time.Sleep(config.StartDelay)
			log.Info("service start delay elapsed")
		}

		restarting := false
		restartInterval := config.InitialRestartInterval

		run := func() error {
			defer func() {
				if err := recover(); err != nil {
					log.Error("service panicked", "msg", err)
					util.PrintStackTrace(3)
					log.Error("restarting service", "restart_interval", restartInterval.String())
					time.Sleep(restartInterval)
					restartInterval = config.NextRestartInterval(restartInterval)
					restarting = true
				}
			}()

			if env.Verbose {
				log.Info("service started")
			}

			restarting = false
			return config.Service.Run()
		}

		for {
			if err := run(); err != nil {
				// If run returns an error then this service has flagged an unrecoverable situation, so
				// do the only thing we can really do here, and blow up.
				panic(err)
			}

			if config.RunInterval < 0 {
				break
			}

			if !restarting {
				time.Sleep(config.RunInterval)
			}
		}
	}()
}

// ExponentialBackoff transforms a time.Duration to the next time.Duration using exponential
// backoff.
func ExponentialBackoff(d time.Duration) time.Duration {
	return d * 2
}
