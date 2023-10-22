package globals

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

type Values struct {
	Flash *flash.Message
	User  *m.User
}

func User(c *fiber.Ctx) *m.User {
	return c.Locals("user").(*m.User)
}

func IsAdmin(c *fiber.Ctx) bool {
	return User(c).Admin
}

func Get(c *fiber.Ctx) *Values {
	values := Values{}

	if flashMessage, found := session.GetOnce[flash.Message](c, "globals.flash"); found {
		values.Flash = &flashMessage
	}

	values.User = User(c)

	return &values
}
