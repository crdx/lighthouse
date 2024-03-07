package device

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	device := util.Param[m.Device](c)

	device.Delete()

	auditLogR.Add(c, "Deleted device %s", device.AuditName())
	flash.Success(c, "Device deleted")
	return c.Redirect("/")
}
