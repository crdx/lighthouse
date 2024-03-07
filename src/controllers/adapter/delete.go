package adapter

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	adapter := util.Param[m.Adapter](c)

	adapter.Delete()

	auditLogR.Add(c, "Deleted adapter %s", adapter.AuditName())
	flash.Success(c, "Adapter deleted")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
