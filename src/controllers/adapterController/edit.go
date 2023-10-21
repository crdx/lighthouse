package adapterController

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
)

type Form struct {
	Name   string `form:"name" validate:"max=100" transform:"trim"`
	Vendor string `form:"vendor" validate:"max=100" transform:"trim"`
}

func get(c *fiber.Ctx) (*m.Device, *m.Adapter, bool) {
	adapter, found := db.First[m.Adapter](c.Params("id"))
	if !found {
		return nil, nil, false
	}

	device, found := adapter.Device()
	if !found {
		return nil, nil, false
	}

	return device, adapter, true
}

func ViewEdit(c *fiber.Ctx) error {
	device, adapter, found := get(c)
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("adapters/edit", fiber.Map{
		"adapter": adapter,
		"device":  device,
		"fields":  validate.Fields[Form](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	device, adapter, found := get(c)

	if !found {
		return c.SendStatus(400)
	}

	form := new(Form)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
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

	flash.Success(c, "Adapter saved")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
