package deviceController

import (
	"fmt"

	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
)

func ViewEdit(c *fiber.Ctx) error {
	device, found := getDevice(c, c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/edit", fiber.Map{
		"device":  device,
		flash.Key: c.Locals(flash.Key),
	})
}

func Edit(c *fiber.Ctx) error {
	device, found := getDevice(c, c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	type Form struct {
		Name        string `form:"name" validate:"required,max=100"`
		Icon        string `form:"icon" validate:"required,max=100"`
		Notes       string `form:"notes" validate:"max=4096"`
		GracePeriod uint   `form:"grace_period" validate:"required,gte=1,lte=60"`
	}

	form := new(Form)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	if fields, err := validate.Struct(form); err {
		return c.Render("edit", fiber.Map{
			"device":  device,
			"err":     err,
			"fields":  fields,
			flash.Key: flash.GetFailure("Unable to save device"),
		})
	}

	device.Update(reflectutil.StructToMap(form, "form"))

	flash.AddSuccess(c, "Device saved")
	return c.Redirect(fmt.Sprintf("/device/%d", device.ID))
}