package audit

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/util/dbutil"
	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "audit",
		"mode":    "list",
		"log":     db.B[m.AuditLog]().Order("created_at DESC").Find(),
		"users":   dbutil.MapByID(db.B[m.User]().Find()),
		"devices": dbutil.MapByID(db.B[m.Device]().Find()),
		"globals": globals.Get(c),
	})
}
