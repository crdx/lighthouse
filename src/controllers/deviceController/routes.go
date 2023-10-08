package deviceController

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/", Get)

	deviceGroup := app.Group("/device/:id<int>")
	deviceGroup.Get("/", View)
	deviceGroup.Post("/delete", Delete)
	deviceGroup.Get("/edit", ViewEdit)
	deviceGroup.Post("/edit", Edit)
	deviceGroup.Post("/merge", Merge)
}
