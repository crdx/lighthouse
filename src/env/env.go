package env

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	ModeDevelopment = "development"
	ModeProduction  = "production"

	LogTypeAll    = "all"
	LogTypeDisk   = "disk"
	LogTypeStderr = "stderr"
	LogTypeNone   = "none"

	AuthTypeBasic = "basic"
	AuthTypeNone  = "none"
)

var (
	Debug = len(os.Getenv("LIGHTHOUSE_DEBUG")) > 0

	Production = os.Getenv("MODE") == ModeProduction
	BindHost   = os.Getenv("HOST")
	BindPort   = os.Getenv("PORT")

	LogType = os.Getenv("LOG_TYPE")
	LogPath = os.Getenv("LOG_PATH")

	AuthType = os.Getenv("AUTH_TYPE")
	AuthUser = os.Getenv("AUTH_USER")
	AuthPass = os.Getenv("AUTH_PASS")

	DatabaseName     = os.Getenv("DB_NAME")
	DatabaseUser     = os.Getenv("DB_USER")
	DatabasePass     = os.Getenv("DB_PASS")
	DatabaseSocket   = os.Getenv("DB_SOCK")
	DatabaseHost     = os.Getenv("DB_HOST")
	DatabaseCharSet  = os.Getenv("DB_CHARSET")
	DatabaseTimeZone = os.Getenv("DB_TZ")

	EnableLiveReload = os.Getenv("LIVE_RELOAD") != ""
)

func Check() {
	// In development no port means use a random port, but this will never be correct for production.
	if Production && BindPort == "" {
		panic("running in production (MODE=production) but no port set")
	}

	require("HOST")

	requireIn("AUTH_TYPE", []string{"none", "basic"}, true)

	if AuthType == AuthTypeBasic {
		require("AUTH_USER")
		require("AUTH_PASS")
	}

	if DatabaseSocket == "" && DatabaseHost == "" {
		panic("required environment variable DB_SOCK or DB_HOST is not set")
	}

	require("DB_NAME")
	require("DB_USER")

	requireIn("LOG_TYPE", []string{"all", "disk", "stderr", "none"}, false)

	if LogType == LogTypeAll || LogType == LogTypeDisk {
		require("LOG_PATH")
	}
}

func require(name string) {
	if os.Getenv(name) == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", name))
	}
}

func requireIn(name string, values []string, canBeEmpty bool) {
	if !canBeEmpty {
		require(name)
	}

	value := os.Getenv(name)

	if canBeEmpty && value == "" {
		return
	}

	if !slices.Contains(values, value) {
		s := ""
		if canBeEmpty {
			s = `, or the empty string ("")`
		}

		panic(fmt.Sprintf(
			`required environment variable %s contains an invalid value (must be one of: "%s"%s)`,
			name,
			strings.Join(values, `", "`),
			s,
		))
	}
}
