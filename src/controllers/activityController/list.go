package activityController

import (
	"fmt"
	"html/template"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceStateLogR"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/pager"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func List(c *fiber.Ctx) error {
	pageNumber, ok := pager.GetCurrentPageNumber(c.QueryInt(pager.Key, 1))

	if !ok {
		return c.SendStatus(404)
	}

	deviceID := uint(c.QueryInt("device_id", 0))

	rows := deviceStateLogR.GetListView(deviceID, pageNumber, constants.ActivityRowsPerPage)
	rowCount := deviceStateLogR.GetListViewRowCount(deviceID)
	pageCount := pager.GetPageCount(rowCount, constants.ActivityRowsPerPage)

	if pageNumber > pageCount {
		return c.SendStatus(404)
	}

	templateParams := fiber.Map{
		"rows":            rows,
		"typeColumnLabel": template.HTML(constants.TypeColumnLabel),
		"globals":         globals.Get(c),
	}

	queryParams := map[string]string{}

	if deviceID != 0 {
		if device, found := db.For[m.Device](deviceID).First(); !found {
			return c.SendStatus(404)
		} else {
			templateParams["device"] = device
		}
		queryParams["device_id"] = fmt.Sprint(deviceID)
	}

	templateParams["pagingState"] = lo.Must(pager.GetState(
		pageNumber,
		pageCount,
		"/activity",
		queryParams,
	))

	return c.Render("activity/list", templateParams)
}
