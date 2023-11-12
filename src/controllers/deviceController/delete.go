package deviceController

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	device := c.Locals("device").(*m.Device)

	device.Delete()

	auditLogR.Add(c, "Deleted device %s", device.AuditName())
	flash.Success(c, "Device deleted")
	return c.Redirect("/")
}
