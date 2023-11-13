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
	deviceGroup.Post("/delete", auth.Editor, Delete)
	deviceGroup.Get("/edit", auth.Editor, ViewEdit)
	deviceGroup.Post("/edit", auth.Editor, Edit)
	deviceGroup.Post("/merge", auth.Editor, Merge)
}
