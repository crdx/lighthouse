package audit

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/userR"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "audit",
		"mode":    "list",
		"log":     db.B[m.AuditLog]().Order("created_at DESC").Find(),
		"users":   userR.Map(),
		"globals": globals.Get(c),
	})
}
