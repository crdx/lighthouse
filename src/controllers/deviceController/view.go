package deviceController

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/deviceStateLogR"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device := c.Locals("device").(*m.Device)

	rows := 6
	if device.Notes != "" {
		rows++
	}

	return c.Render("devices/view", fiber.Map{
		"mode":     "view",
		"device":   device,
		"devices":  deviceR.All(),
		"adapters": device.Adapters(),
		"activity": deviceStateLogR.LatestActivityForDevice(device.ID, rows),
		"globals":  globals.Get(c),
	})
}
