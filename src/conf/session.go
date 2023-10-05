package conf

import (
	"crdx.org/session"
)

func GetSessionConfig() *session.Config {
	return &session.Config{
		Table: "sessions",

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
