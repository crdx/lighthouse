package deviceController

import (
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/repos/deviceR"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device, found := getDevice(c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/view", fiber.Map{
		"device":   device,
		"adapters": device.Adapters(),
		flash.Key:  c.Locals(flash.Key),
	})
}
