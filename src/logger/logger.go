package logger

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"crdx.org/lighthouse/env"
	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
)

var logger *slog.Logger
var mutex sync.Mutex

func New() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()

	if logger != nil {
		return logger
	}

	switch env.LogType {
	case env.LogTypeAll:
		logger = slog.New(slogmulti.Fanout(getDiskHandler(), getStderrHandler()))
	case env.LogTypeDisk:
		logger = slog.New(getDiskHandler())
	case env.LogTypeStdout:
		logger = slog.New(getStderrHandler())
	default:
		panic("unexpected env.LogType")
	}

	logger.Info("logger init complete", "type", env.LogType)
	return logger
}

func getDiskHandler() slog.Handler {
	file := lo.Must(os.OpenFile("logs/lighthouse.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644))
	lo.Must0(os.MkdirAll("logs", 0755))
	return slog.NewJSONHandler(file, nil)
}

func getStderrHandler() slog.Handler {
	return tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Kitchen,
	})
}
