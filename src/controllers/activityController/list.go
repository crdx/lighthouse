package activityController

import (
	"fmt"
	"html/template"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/pager"
	"crdx.org/lighthouse/repos/deviceStateLogR"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func List(c *fiber.Ctx) error {
	pageNumber, ok := pager.GetCurrentPageNumber(c)

	if !ok {
		return c.SendStatus(404)
	}

	deviceID := uint(c.QueryInt("device_id", 0))

	entries := deviceStateLogR.GetListView(deviceID, pageNumber, constants.ActivityEntriesPerPage)
	totalEntries := deviceStateLogR.GetListViewTotal(deviceID)
	pageCount := pager.GetPageCount(totalEntries, constants.ActivityEntriesPerPage)

	if pageCount > 0 && pageNumber > pageCount {
		return c.SendStatus(404)
	}

	templateParams := fiber.Map{
		"entries":         entries,
		"typeColumnLabel": template.HTML(constants.TypeColumnLabel),
		flash.Key:         c.Locals(flash.Key),
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
