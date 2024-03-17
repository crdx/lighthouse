package device

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	device := parseparam.Get[db.Device](c)

	device.CascadeDelete()

	auditLogR.Add(c, "Deleted device %s", device.AuditName())
	flash.Success(c, "Device deleted")
	return c.Redirect("/")
}
