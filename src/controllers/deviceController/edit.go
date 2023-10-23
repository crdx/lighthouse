package deviceController

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

type EditForm struct {
	Name        string `form:"name" validate:"max=100" transform:"trim"`
	Icon        string `form:"icon" validate:"max=100" transform:"trim"`
	Notes       string `form:"notes" validate:"max=5000" transform:"trim"`
	GracePeriod string `form:"grace_period" validate:"required,number" transform:"trim"`
	Watch       bool   `form:"watch"`
}

func ViewEdit(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}

	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/edit", fiber.Map{
		"mode":    "edit",
		"fields":  validate.Fields[EditForm](),
		"device":  device,
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}

	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	form := new(EditForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

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
