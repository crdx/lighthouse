package config

import (
	"time"

	"crdx.org/lighthouse/pkg/session"
	"github.com/gofiber/fiber/v3"
)

const tableName = "sessions"

// NewSession creates session middleware with the production config.
func NewSession(dsn string) fiber.Handler {
	return session.New(&session.Config{
		DSN:          dsn,
		Table:        tableName,
		IdleTimeout:  365 * 24 * time.Hour,
		CookieSecure: false,
	})
}

// NewTestSession creates session middleware with the test config.
func NewTestSession(dsn string) fiber.Handler {
	return session.New(&session.Config{
		DSN:          dsn,
		Table:        tableName,
		CookieSecure: false,
	})
}
