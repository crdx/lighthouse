package deviceController

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/deviceStateLogR"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("devices/view", fiber.Map{
		"mode":     "view",
		"device":   device,
		"devices":  deviceR.All(),
		"adapters": device.Adapters(),
		"activity": deviceStateLogR.LatestActivityForDevice(device.ID),
		"globals":  globals.Get(c),
	})
}
