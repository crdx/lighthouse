package deviceController

import (
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func DeleteDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	device.Delete()
	flash.AddSuccess(c, "Device deleted")
	return c.Redirect("/")
}
