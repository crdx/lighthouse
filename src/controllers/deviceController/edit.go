package deviceController

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	Name        string `form:"name" validate:"max=100" transform:"trim"`
	Icon        string `form:"icon" validate:"max=100" transform:"trim"`
	Notes       string `form:"notes" validate:"max=5000" transform:"trim"`
	GracePeriod string `form:"grace_period" validate:"required,number,grace_period" transform:"trim"`
	Watch       bool   `form:"watch"`
	Limit       string `form:"limit" validate:"required,number" transform:"trim"`
}

func ViewEdit(c *fiber.Ctx) error {
	device := c.Locals("device").(*m.Device)

	return c.Render("devices/edit", fiber.Map{
		"mode":    "edit",
		"fields":  validate.Fields[EditForm](),
		"device":  device,
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	device := c.Locals("device").(*m.Device)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
		flash.Failure(c, "Unable to save device")

		return c.Render("devices/edit", fiber.Map{
			"device":  device,
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	device.Update(reflectutil.StructToMap(form, "form"))

	flash.Success(c, "Device saved")
	return c.Redirect(fmt.Sprintf("/device/%d", device.ID))
}
