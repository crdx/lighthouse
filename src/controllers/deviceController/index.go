package deviceController

import (
	"slices"

	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/tpl"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
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

	return c.Render("devices/index", fiber.Map{
		"devices": deviceR.GetListView(currentSortColumn, currentSortDirection),
		"columns": tpl.AddSortMetadata(currentSortColumn, currentSortDirection, columns),
		flash.Key: c.Locals(flash.Key),
	})
}
