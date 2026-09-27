package activity

import (
	"fmt"
	"html/template"
	"maps"
	"slices"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceStateLogR"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/pager"
	"crdx.org/lighthouse/pkg/util/tplutil"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func List(c fiber.Ctx) error {
	pageNumber, ok := pager.GetCurrentPageNumber(fiber.Query[int](c, pager.Key, 1))

	if !ok {
		return c.SendStatus(404)
	}

	typeColumnLabel := template.HTML(constants.TypeColumnLabel) //nolint:gosec // Constant string, not user input.
	columns := map[string]tplutil.ColumnConfig{
		"type":  {Label: typeColumnLabel, DefaultSortDirection: "asc", Minimal: true},
		"name":  {Label: "Name", DefaultSortDirection: "asc"},
		"state": {Label: "State", DefaultSortDirection: "asc"},
		"when":  {Label: "When", DefaultSortDirection: "desc"},
		"date":  {Label: "Date", DefaultSortDirection: "desc"},
	}
	currentSortColumn := c.Query("sc", "when")
	currentSortDirection := c.Query("sd", "desc")

	if !slices.Contains([]string{"asc", "desc"}, currentSortDirection) {
		return c.SendStatus(400)
	}

	if !slices.Contains(slices.Collect(maps.Keys(columns)), currentSortColumn) {
		return c.SendStatus(400)
	}

	deviceID := int64(fiber.Query[int](c, "device_id", 0))
	queryParams := map[string]string{
		"sc": currentSortColumn,
		"sd": currentSortDirection,
	}
	headingQueryParams := map[string]string{}

	rows := deviceStateLogR.GetList(
		pageNumber,
		constants.ActivityRowsPerPage,
		deviceID,
		currentSortColumn,
		currentSortDirection,
	)
	rowCount := deviceStateLogR.GetListRowCount(deviceID)
	pageCount := pager.GetPageCount(rowCount, constants.ActivityRowsPerPage)

	if pageNumber > pageCount {
		return c.SendStatus(404)
	}

	templateParams := fiber.Map{
		"rows":    rows,
		"globals": globals.Get(c),
	}

	if deviceID != 0 {
		if device, found := db.FindDevice(deviceID); !found {
			return c.SendStatus(404)
		} else {
			templateParams["device"] = device
		}
		queryParams["device_id"] = fmt.Sprint(deviceID)
		headingQueryParams["device_id"] = fmt.Sprint(deviceID)
	}

	templateParams["columns"] = tplutil.AddMetadata(
		currentSortColumn,
		currentSortDirection,
		headingQueryParams,
		columns,
	)
	templateParams["pagingState"] = lo.Must(pager.GetState(
		pageNumber,
		pageCount,
		"/activity",
		queryParams,
	))

	return c.Render("activity/list", templateParams)
}
