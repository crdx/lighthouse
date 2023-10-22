package conf

import (
	"time"

	"crdx.org/session"
)

func GetSessionConfig() *session.Config {
	return &session.Config{
		Table:      "sessions",
		Expiration: 365 * 24 * time.Hour,

		// lighthouse is intended to be accessed over the local network only.
		CookieSecure: false,
	}
}

func GetTestSessionConfig() *session.Config {
	return &session.Config{
		Table:        "sessions",
		CookieSecure: false,
	}
}
