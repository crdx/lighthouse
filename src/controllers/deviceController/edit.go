package deviceController

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
)

func ViewEdit(c *fiber.Ctx) error {
	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/edit", fiber.Map{
		"device":  device,
		flash.Key: c.Locals(flash.Key),
	})
}

func Edit(c *fiber.Ctx) error {
	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	type Form struct {
		Name        string `form:"name" validate:"required,max=100" transform:"trim"`
		Icon        string `form:"icon" validate:"required,max=100" transform:"trim"`
		Notes       string `form:"notes" validate:"max=4096" transform:"trim"`
		GracePeriod string `form:"grace_period" validate:"required,number,gte=1,lte=60" transform:"trim"`
		Watch       bool   `form:"watch"`
	}

	form := new(Form)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
		return c.Render("devices/edit", fiber.Map{
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
