package deviceController

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"github.com/gofiber/fiber/v2"
)

func Merge(c *fiber.Ctx) error {
	device1, found := getDevice(c.Params("id"))

	if !found {
		return c.SendStatus(400)
	}

	device2, found := getDevice(c.FormValue("device_id"))

	if !found {
		return c.SendStatus(400)
	}

	// The parent device will be the one that was discovered first.
	var parent, child *m.Device
	if device1.CreatedAt.Before(device2.CreatedAt) {
		parent, child = device1, device2
	} else {
		parent, child = device2, device1
	}

	for _, adapter := range child.Adapters() {
		adapter.Update("device_id", parent.ID)
	}

	if child.LastSeen.After(parent.LastSeen) {
		parent.Update("last_seen", child.LastSeen)
	}

	child.Delete()

	flash.AddSuccess(c, fmt.Sprintf(
		"Device %s merged into %s (as latter was discovered first)",
		child.Name,
		parent.Name,
	))

	return c.Redirect(fmt.Sprintf("/device/%d", parent.ID))
}