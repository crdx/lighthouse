package audit

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v3"
)

func List(c fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "audit",
		"mode":    "list",
		"log":     db.FindAuditLogsSorted(),
		"users":   db.MapByID(db.FindUsersUnscoped()),
		"devices": db.MapByID(db.FindDevicesUnscoped()),
		"globals": globals.Get(c),
	})
}
