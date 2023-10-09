package deviceController

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/repos/deviceR"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/view", fiber.Map{
		"device":   device,
		"devices":  deviceR.All(),
		"adapters": device.Adapters(),
		flash.Key:  c.Locals(flash.Key),
	})
}
