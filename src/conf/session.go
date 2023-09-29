package conf

import (
	"crdx.org/lighthouse/env"
	"crdx.org/session"
)

func GetSessionConfig() *session.Config {
	return &session.Config{
		Table:        "sessions",
		CookieSecure: env.Production,
	}
}

func GetTestSessionConfig() *session.Config {
	return &session.Config{
		Table:        "sessions",
		CookieSecure: false,
	}
}
