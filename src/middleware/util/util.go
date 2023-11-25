package util

import (
	"strings"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"github.com/gofiber/fiber/v2"
)

const parseParamPrefix = "parseParam_"

// NewParseParam returns middleware that looks for a route parameter corresponding to an ID,
// instantiates the model T with it, and assigns it to c.Locals(name).
func NewParseParam[T any](param string, name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance, found := db.First[T](c.Params(param))
		if !found {
			return c.SendStatus(404)
		}
		c.Locals(parseParamPrefix+name, instance)
		return c.Next()
	}
}

func Param[T any](c *fiber.Ctx) *T {
	var v T
	name := strings.ToLower(reflectutil.GetType(v).Name())
	return c.Locals(parseParamPrefix + name).(*T)
}
