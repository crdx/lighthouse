package adapterController

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	adapter, found := db.First[m.Adapter](c.Params("id"))
	if !found {
		return c.SendStatus(400)
	}

	adapter.Delete()
	flash.AddSuccess(c, "Adapter deleted")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}