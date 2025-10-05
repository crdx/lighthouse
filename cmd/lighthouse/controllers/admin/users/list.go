package users

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v3"
)

func List(c fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "list",
		"users":   db.FindUsersSorted(),
		"globals": globals.Get(c),
	})
}
