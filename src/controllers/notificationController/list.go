package notificationController

import (
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m/repo/notificationR"
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

	rows := notificationR.GetListView(pageNumber, constants.ActivityRowsPerPage)
	rowCount := notificationR.GetListViewRowCount()
	pageCount := pager.GetPageCount(rowCount, constants.NotificationRowsPerPage)

	if pageNumber > pageCount {
		return c.SendStatus(404)
	}

	templateParams := fiber.Map{
		"rows":    rows,
		"globals": globals.Get(c),
	}

	templateParams["pagingState"] = lo.Must(pager.GetState(
		pageNumber,
		pageCount,
		"/notifications",
		nil,
	))

	return c.Render("notifications/list", templateParams)
}
