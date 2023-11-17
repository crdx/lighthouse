package util

import (
	"crdx.org/db"
	"github.com/gofiber/fiber/v2"
)

// NewParseParam returns middleware that looks for a route parameter corresponding to an ID,
// instantiates the model T with it, and assigns it to c.Locals(name).
func NewParseParam[T any](param string, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance, found := db.First[T](c.Params(param))
		if !found {
			return c.SendStatus(404)
		}
		c.Locals(name, instance)
		return c.Next()
	}
}
