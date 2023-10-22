package adminController

import (
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}

	return c.Redirect("/admin/settings")
}
