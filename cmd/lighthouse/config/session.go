package config

import (
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/session"
	"github.com/gofiber/fiber/v3"
)

// NewSessionMiddleware creates session middleware with the production config.
func NewSessionMiddleware(dbConfig *db.Config) fiber.Handler {
	return session.New(&session.Config{
		Table:       "sessions",
		IdleTimeout: 365 * 24 * time.Hour,

		// lighthouse is intended to be accessed over the local network only.
		CookieSecure: false,
	}, dbConfig.DataSource.Format())
}

// NewTestSessionMiddleware creates session middleware with the test config.
func NewTestSessionMiddleware(dbConfig *db.Config) fiber.Handler {
	return session.New(&session.Config{
		Table:        "sessions",
		CookieSecure: false,
	}, dbConfig.DataSource.Format())
}
