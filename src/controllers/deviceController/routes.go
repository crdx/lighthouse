package deviceController

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/util"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/", List)

	deviceGroup := app.Group("/device/:id<int>").
		Use(util.NewParseParam[m.Device]("id", "device"))

	deviceGroup.Get("/", View)
	deviceGroup.Post("/delete", auth.Admin, Delete)
	deviceGroup.Get("/edit", auth.Admin, ViewEdit)
	deviceGroup.Post("/edit", auth.Admin, Edit)
	deviceGroup.Post("/merge", auth.Admin, Merge)
}
