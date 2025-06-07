package logger

import (
	"log/slog"
	"os"
	"sync"

	"crdx.org/lighthouse/pkg/env"
	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *slog.Logger
	mutex  sync.Mutex
)

func Get() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()

	if logger != nil {
		return logger
	}

	switch env.LogType() {
	case env.LogTypeAll:
		logger = slog.New(slogmulti.Fanout(getDiskHandler(), getStderrHandler()))
	case env.LogTypeDisk:
		logger = slog.New(getDiskHandler())
	case env.LogTypeStderr:
		logger = slog.New(getStderrHandler())
	case env.LogTypeNone:
		logger = slog.New(slog.DiscardHandler)
	}

	logger.Info("logger init complete", "type", env.LogType())
	return logger
}

func With(args ...any) *slog.Logger {
	return Get().With(args...)
}

func getDiskHandler() slog.Handler {
	return slog.NewJSONHandler(&lumberjack.Logger{
		Filename: env.LogPath(),
		MaxSize:  50, // MB
		Compress: true,
	}, nil)
}

func getStderrHandler() slog.Handler {
	return tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: "15:04:05 |", // Closer in style to fiber's debug output.
	})
}
