package users

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "list",
		"users":   db.B[m.User]().Order("username ASC").Find(),
		"globals": globals.Get(c),
	})
}
