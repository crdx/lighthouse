package service

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	Name string `form:"name" validate:"max=100" transform:"trim"`
}

func ViewEdit(c *fiber.Ctx) error {
	service := util.Param[m.Service](c)

	return c.Render("services/edit", fiber.Map{
		"service": service,
		"device":  service.Device(),
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	service := util.Param[m.Service](c)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
		flash.Failure(c, "Unable to save service")

		return c.Render("services/edit", fiber.Map{
			"service": service,
			"device":  service.Device(),
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	service.Update(reflectutil.StructToMap(form, "form"))

	auditLogR.Add(c, "Edited service %s", service.Fresh().AuditName())
	flash.Success(c, "Service saved")
	return c.Redirect(fmt.Sprintf("/device/%d", service.DeviceID))
}
