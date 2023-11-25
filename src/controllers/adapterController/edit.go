package adapterController

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
	Name   string `form:"name" validate:"max=100" transform:"trim"`
	Vendor string `form:"vendor" validate:"max=100" transform:"trim"`
}

func ViewEdit(c *fiber.Ctx) error {
	adapter := util.Param[m.Adapter](c)
	device := adapter.Device()

	return c.Render("adapters/edit", fiber.Map{
		"adapter": adapter,
		"device":  device,
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	adapter := util.Param[m.Adapter](c)
	device := adapter.Device()

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
		flash.Failure(c, "Unable to save adapter")

		return c.Render("adapters/edit", fiber.Map{
			"adapter": adapter,
			"device":  device,
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	adapter.Update(reflectutil.StructToMap(form, "form"))

	auditLogR.Add(c, "Edited adapter %s", adapter.Fresh().AuditName())
	flash.Success(c, "Adapter saved")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
