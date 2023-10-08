package deviceController

import (
	"crdx.org/lighthouse/m"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getDevice(c *fiber.Ctx) (*m.Device, bool) {
	return m.ForDevice(uint(lo.Must(c.ParamsInt("id")))).First()
}
