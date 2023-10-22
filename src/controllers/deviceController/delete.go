package deviceController

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}

	device, found := db.First[m.Device](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	device.Delete()
	flash.Success(c, "Device deleted")
	return c.Redirect("/")
}
