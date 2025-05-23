package device

import (
	"html/template"
	"maps"
	"slices"

	"crdx.org/lighthouse/db/repo/deviceR"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/util/tplutil"
	"github.com/gofiber/fiber/v2"
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

	currentSortColumn := c.Query("sc", "ip")
	currentSortDirection := c.Query("sd", "asc")
	currentFilter := c.Query("f", "online")

	if !slices.Contains([]string{"asc", "desc"}, currentSortDirection) {
		return c.SendStatus(400)
	}

	if !slices.Contains(slices.Collect(maps.Keys(columns)), currentSortColumn) {
		return c.SendStatus(400)
	}

	return c.Render("device/list", fiber.Map{
		"currentFilter": currentFilter,
		"devices":       deviceR.GetList(currentSortColumn, currentSortDirection, currentFilter),
		"counts":        deviceR.GetListCounts(),
		"columns":       tplutil.AddMetadata(currentSortColumn, currentSortDirection, currentFilter, columns),
		"globals":       globals.Get(c),
	})
}
