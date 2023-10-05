package indexController

import (
	"slices"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/tpl"
	"golang.org/x/exp/maps"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func Get(c *fiber.Ctx) error {
	columns := map[string]tpl.SortableColumnConfig{
		"name":   {Label: "Name", DefaultSortDirection: "asc"},
		"ip":     {Label: "IP Address", DefaultSortDirection: "asc"},
		"vendor": {Label: "Vendor", DefaultSortDirection: "asc"},
		"mac":    {Label: "MAC Address", DefaultSortDirection: "asc"},
		"seen":   {Label: "Last Seen", DefaultSortDirection: "desc"},
	}

	currentSortColumn := c.Query("sc", "seen")
	currentSortDirection := c.Query("sd", "desc")

	if !slices.Contains([]string{"asc", "desc"}, currentSortDirection) {
		return c.SendStatus(400)
	}

	if !slices.Contains(maps.Keys(columns), currentSortColumn) {
		return c.SendStatus(400)
	}

	return c.Render("index", fiber.Map{
		"devices": deviceR.GetListView(currentSortColumn, currentSortDirection),
		"columns": tpl.AddSortMetadata(currentSortColumn, currentSortDirection, columns),
		flash.Key: c.Locals(flash.Key), // TODO: Find a way to not need this in every route.
	})
}

func ViewDevice(c *fiber.Ctx) error {
	if device, found := m.ForDevice(uint(lo.Must(c.ParamsInt("id")))).First(); found {
		return c.Render("view", fiber.Map{
			"device":   device,
			"adapters": device.Adapters(),
			flash.Key:  c.Locals(flash.Key), // TODO: Find a way to not need this in every route.
		})
	}

	return c.SendStatus(404)
}

func DeleteDevice(c *fiber.Ctx) error {
	if device, found := m.ForDevice(uint(lo.Must(c.ParamsInt("id")))).First(); found {
		device.Delete()
		flash.AddSuccess(c, "Device deleted")
	}

	return c.Redirect("/")
}
