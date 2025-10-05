package device

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	app.Get("/", List).Name("devices")

	deviceGroup := app.Group("/device/:id<int>").
		Use(parseparam.New("id", db.FindDevice))

	deviceGroup.Get("/", View).Name("devices")
	deviceGroup.Post("/delete", auth.Editor, Delete)
	deviceGroup.Get("/edit", auth.Editor, ViewEdit).Name("devices")
	deviceGroup.Post("/edit", auth.Editor, Edit)
	deviceGroup.Post("/merge", auth.Editor, Merge)
}
