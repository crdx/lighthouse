package service

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
	// Init initialises the service and returns an error if initialisation failed. If an error is
	// returned then the service will not run.
	Init(*Args) error

	// Run runs the service and returns an error for an unrecoverable failure.
	Run() error

	// Config returns the configuration for this service.
	Config() Config
}

type Config interface {
	// Name returns the name of the service.
	Name() string

	// NextRestartInterval returns a function that receives the last restart interval and returns
	// the next one.
	NextRestartInterval(time.Duration) time.Duration

	// InitialRestartInterval returns the initial restart interval.
	InitialRestartInterval() time.Duration

	// RunInterval returns the interval between runs. Zero means no delay and -1 means only run once.
	RunInterval() time.Duration

	// StartDelay returns the initial delay before starting the service.
	InitialStartDelay() time.Duration
}

func Start(service Service) {
	config := service.Config()

	log := logger.New().With("service", config.Name())

	if err := service.Init(&Args{Logger: log}); err != nil {
		log.Error("service initialisation failed", "msg", err)
		return
	}

	restartInterval := config.InitialRestartInterval()
	runInterval := config.RunInterval()
	initialStartDelay := config.InitialStartDelay()

	log.Info(
		"service initialisation complete",
		"run_interval", runInterval.String(),
		"restart_interval", restartInterval.String(),
		"initial_delay", initialStartDelay.String(),
	)

	if initialStartDelay > 0 {
		util.Sleep(initialStartDelay)
		log.Info("service start delay elapsed")
	}

	restarting := false

	run := func() error {
		defer func() {
			if err := recover(); err != nil {
				log.Error("service panicked", "msg", err)
				util.PrintStackTrace(3)
				log.Error("restarting service", "restart_interval", restartInterval.String())
				util.Sleep(restartInterval)
				restartInterval = config.NextRestartInterval(restartInterval)
				restarting = true
			}
		}()

		if env.Verbose {
			log.Info("service started")
		}

		restarting = false
		return service.Run()
	}

	for {
		if err := run(); err != nil {
			// If run returns an error then this service has flagged an unrecoverable situation, so
			// do the only thing we can really do here, and blow up.
			panic(err)
		}

		if runInterval < 0 {
			break
		}

		if !restarting {
			util.Sleep(runInterval)
		}
	}
}

// ExponentialBackoff transforms a time.Duration to the next time.Duration using exponential
// backoff.
func ExponentialBackoff(d time.Duration) time.Duration {
	return d * 2
}
