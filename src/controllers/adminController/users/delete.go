package users

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	user, found := db.First[m.User](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	// Current user can't delete themselves.
	if globals.User(c).ID == user.ID {
		return c.SendStatus(400)
	}

	user.Delete()

	flash.Success(c, "User deleted")
	return c.Redirect("/admin/users")
}
