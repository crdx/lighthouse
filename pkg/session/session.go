package session

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/samber/lo"
)

type Config struct {
	Table        string        // The table to store the session data in.
	CookieSecure bool          // Whether the cookie should be HTTPS-only.
	IdleTimeout  time.Duration // How long the session cookie should last.
}

// New initialises the session and returns a middleware handler.
func New(config *Config, dsn string) fiber.Handler {
	return session.New(session.Config{
		Storage: mysql.New(mysql.Config{
			ConnectionURI: dsn,
			Table:         config.Table,
		}),
		Extractor:      extractors.FromCookie("session"),
		CookieSecure:   config.CookieSecure,
		CookieHTTPOnly: true,
		IdleTimeout:    config.IdleTimeout,
	})
}

// Get fetches a value from the session as T. If the session doesn't contain a T then the zero value
// of T is returned.
func Get[T any](c fiber.Ctx, key string) T {
	value, _ := TryGet[T](c, key)
	return value
}

// GetOnce fetches a value from the session as T, then deletes it. If the session doesn't contain a
// T then the zero value of T is returned.
func GetOnce[T any](c fiber.Ctx, key string) T {
	value, _ := TryGetOnce[T](c, key)
	return value
}

// TryGet fetches a value from the session as T, returning the value and a boolean indicating
// whether the value was present in the session.
func TryGet[T any](c fiber.Ctx, key string) (T, bool) {
	if value, found := get(c, key).(T); found {
		return value, true
	}
	var value T
	return value, false
}

// TryGetOnce fetches a value from the session as T, then deletes it from the session.
func TryGetOnce[T any](c fiber.Ctx, key string) (T, bool) {
	if value, found := getOnce(c, key).(T); found {
		return value, true
	}
	var value T
	return value, false
}

// Set stores a value in the session. If T is a custom type then it may need to be registered with
// gob.Register first.
func Set[T any](c fiber.Ctx, key string, value T) {
	session.FromContext(c).Set(key, value)
}

// Delete deletes a value from the session.
func Delete(c fiber.Ctx, key string) {
	session.FromContext(c).Delete(key)
}

// Destroy destroys the session.
func Destroy(c fiber.Ctx) {
	lo.Must0(session.FromContext(c).Destroy())
}

// GetID returns the session ID.
func GetID(c fiber.Ctx) string {
	return session.FromContext(c).ID()
}

func get(c fiber.Ctx, key string) any {
	return session.FromContext(c).Get(key)
}

func getOnce(c fiber.Ctx, key string) any {
	value := session.FromContext(c).Get(key)
	session.FromContext(c).Delete(key)
	return value
}
