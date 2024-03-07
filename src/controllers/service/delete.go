package service

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	service := util.Param[m.Service](c)

	service.Delete()

	auditLogR.Add(c, "Deleted service %s", service.AuditName())
	flash.Success(c, "Service deleted")
	return c.Redirect(fmt.Sprintf("/device/%d", service.DeviceID))
}
