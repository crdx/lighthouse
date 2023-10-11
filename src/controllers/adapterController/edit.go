package adapterController

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func ViewEdit(c *fiber.Ctx) error {
	adapter, found := db.First[m.Adapter](c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("adapters/edit", fiber.Map{
		"adapter": adapter,
		"device":  lo.Must(adapter.Device()),
		flash.Key: c.Locals(flash.Key),
	})
}

func Edit(c *fiber.Ctx) error {
	adapter, found := db.First[m.Adapter](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	type Form struct {
		Name   string `form:"name" validate:"max=100" transform:"trim"`
		Vendor string `form:"vendor" validate:"max=100" transform:"trim"`
	}

	form := new(Form)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
		return c.Render("adapters/edit", fiber.Map{
			"adapter": adapter,
			"err":     err,
			"fields":  fields,
			flash.Key: flash.GetFailure("Unable to save adapter"),
		})
	}

	adapter.Update(reflectutil.StructToMap(form, "form"))

	flash.AddSuccess(c, "Adapter saved")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
