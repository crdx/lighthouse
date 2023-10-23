package adapterController

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/util"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	adapterGroup := app.Group("/adapter/:id<int>").
		Use(auth.Admin).
		Use(util.NewParseParam[m.Adapter]("id", "adapter"))

	adapterGroup.Post("/delete", Delete)
	adapterGroup.Get("/edit", ViewEdit)
	adapterGroup.Post("/edit", Edit)
}
