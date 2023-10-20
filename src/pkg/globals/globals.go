package globals

import (
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

type Values struct {
	Flash *flash.Message
}

func Get(c *fiber.Ctx) *Values {
	values := Values{}

	if flashMessage, found := session.GetOnce[flash.Message](c, "globals.flash"); found {
		values.Flash = &flashMessage
	}

	return &values
}
