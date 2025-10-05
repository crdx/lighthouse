package service

import (
	"fmt"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v3"
)

func Delete(c fiber.Ctx) error {
	service := parseparam.Get[db.Service](c)

	service.Delete()

	auditLogR.Add(c, "Deleted service %s", service.AuditName())
	flash.Success(c, "Service deleted")
	return c.Redirect().To(fmt.Sprintf("/device/%d", service.DeviceID))
}
