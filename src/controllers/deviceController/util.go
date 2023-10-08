package deviceController

import (
	"strconv"

	"crdx.org/lighthouse/m"
	"github.com/gofiber/fiber/v2"
)

func getDevice(c *fiber.Ctx, v string) (*m.Device, bool) {
	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, false
	}

	return m.ForDevice(uint(i)).First()
}
