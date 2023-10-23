package users

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*m.User)

	// Current user can't delete themselves.
	if globals.User(c).ID == user.ID {
		return c.SendStatus(400)
	}

	user.Delete()

	flash.Success(c, "User deleted")
	return c.Redirect("/admin/users")
}
