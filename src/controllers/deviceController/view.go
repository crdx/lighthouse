package deviceController

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/m/repo/deviceStateLogR"
	"crdx.org/lighthouse/m/repo/probeR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device := util.Param[m.Device](c)

	rows := 6
	if device.Notes != "" {
		rows++
	}

	return c.Render("devices/view", fiber.Map{
		"mode":               "view",
		"device":             device,
		"devices":            deviceR.All(),
		"adapters":           device.Adapters(),
		"services":           device.Services(),
		"serviceTTL":         probeR.TTL(),
		"serviceScanEnabled": settingR.EnableServiceScan(),
		"activity":           deviceStateLogR.LatestActivityForDevice(device.ID, rows),
		"globals":            globals.Get(c),
	})
}
