package adminController

import (
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}

	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"globals": globals.Get(c),
	})
}
