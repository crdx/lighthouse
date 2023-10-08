package deviceController

import (
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	device, found := getDevice(c, c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	device.Delete()
	flash.AddSuccess(c, "Device deleted")
	return c.Redirect("/")
}
