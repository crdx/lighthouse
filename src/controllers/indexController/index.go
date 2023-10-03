package indexController

import (
	"slices"

	"crdx.org/lighthouse/models/deviceM"
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
		"devices": deviceM.GetListView(currentSortColumn, currentSortDirection),
		"columns": tpl.AddSortMetadata(currentSortColumn, currentSortDirection, columns),
	})
}

func ViewDevice(c *fiber.Ctx) error {
	if device, found := deviceM.For(uint(lo.Must(c.ParamsInt("id")))).First(); found {
		return c.Render("view", fiber.Map{
			"device":   device,
			"adapters": device.Adapters(),
		})
	}

	return c.SendStatus(404)
}
