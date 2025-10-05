package activity

import (
	"fmt"
	"html/template"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceStateLogR"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/pager"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func List(c fiber.Ctx) error {
	pageNumber, ok := pager.GetCurrentPageNumber(fiber.Query[int](c, pager.Key, 1))

	if !ok {
		return c.SendStatus(404)
	}

	deviceID := int64(fiber.Query[int](c, "device_id", 0))

	rows := deviceStateLogR.GetList(pageNumber, constants.ActivityRowsPerPage, deviceID)
	rowCount := deviceStateLogR.GetListRowCount(deviceID)
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
		if device, found := db.FindDevice(deviceID); !found {
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
