package deviceController

import (
	"html/template"
	"slices"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/util/tplutil"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
)

func List(c *fiber.Ctx) error {
	watchLabel := template.HTML(constants.WatchColumnLabel)
	typeLabel := template.HTML(constants.TypeColumnLabel)

	columns := map[string]tplutil.ColumnConfig{
		"name":   {Label: "Name", DefaultSortDirection: "asc"},
		"ip":     {Label: "IP Address", DefaultSortDirection: "asc"},
		"vendor": {Label: "Vendor", DefaultSortDirection: "asc"},
		"mac":    {Label: "MAC Address", DefaultSortDirection: "asc"},
		"seen":   {Label: "Last Seen", DefaultSortDirection: "desc"},
		"watch":  {Label: watchLabel, DefaultSortDirection: "desc", Minimal: true},
		"type":   {Label: typeLabel, DefaultSortDirection: "asc", Minimal: true},
	}

	currentSortColumn := c.Query("sc", "seen")
	currentSortDirection := c.Query("sd", "desc")
	currentFilter := c.Query("f", "")

	if !slices.Contains([]string{"asc", "desc"}, currentSortDirection) {
		return c.SendStatus(400)
	}

	if !slices.Contains(maps.Keys(columns), currentSortColumn) {
		return c.SendStatus(400)
	}

	return c.Render("devices/list", fiber.Map{
		"currentFilter": currentFilter,
		"devices":       deviceR.GetListView(currentSortColumn, currentSortDirection, currentFilter),
		"counts":        deviceR.GetCounts(),
		"columns":       tplutil.AddMetadata(currentSortColumn, currentSortDirection, currentFilter, columns),
		"globals":       globals.Get(c),
	})
}
