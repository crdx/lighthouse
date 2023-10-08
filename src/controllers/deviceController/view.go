package deviceController

import (
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func ViewDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/view", fiber.Map{
		"device":   device,
		"adapters": device.Adapters(),
		flash.Key:  c.Locals(flash.Key),
	})
}
