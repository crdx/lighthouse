package parseparam

import (
	"strings"

	"crdx.org/lighthouse/pkg/util/reflectutil"
	"github.com/gofiber/fiber/v3"
)

const prefix = "parseParam_"

// New returns middleware that looks for a route parameter corresponding to an ID, fetches
// the model for it, and assigns it to c.Locals(name).
func New[T any](param string, fetch func(int64) (*T, bool)) fiber.Handler {
	var v T
	name := strings.ToLower(reflectutil.GetType(v).Name())

	return func(c fiber.Ctx) error {
		id := fiber.Params[int](c, param)
		if id == 0 {
			return c.SendStatus(404)
		}

		instance, found := fetch(int64(id))
		if !found {
			return c.SendStatus(404)
		}

		c.Locals(prefix+name, instance)
		return c.Next()
	}
}

// Get returns the route parameter for T.
func Get[T any](c fiber.Ctx) *T {
	var v T
	name := strings.ToLower(reflectutil.GetType(v).Name())
	return c.Locals(prefix + name).(*T)
}
