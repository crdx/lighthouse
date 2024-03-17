package device

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/probeR"
	"crdx.org/lighthouse/db/repo/settingR"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	device := parseparam.Get[db.Device](c)

	var rows int64 = 7
	if device.Notes != "" {
		rows++
	}

	return c.Render("device/view", fiber.Map{
		"mode":               "view",
		"device":             device,
		"devices":            db.FindDevicesSorted(),
		"adapters":           db.FindAdaptersForDevice(device.ID),
		"services":           db.FindServicesForDevice(device.ID),
		"serviceTTL":         probeR.TTL(),
		"serviceScanEnabled": settingR.EnableServiceScan(),
		"activity":           db.FindLatestActivityForDevice(device.ID, rows),
		"globals":            globals.Get(c),
	})
}
