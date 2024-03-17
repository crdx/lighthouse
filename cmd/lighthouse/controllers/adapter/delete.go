package adapter

import (
	"fmt"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	adapter := parseparam.Get[db.Adapter](c)

	adapter.Delete()

	auditLogR.Add(c, "Deleted adapter %s", adapter.AuditName())
	flash.Success(c, "Adapter deleted")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
