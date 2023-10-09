package deviceController

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	device.Delete()
	flash.AddSuccess(c, "Device deleted")
	return c.Redirect("/")
}
