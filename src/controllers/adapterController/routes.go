package adapterController

import (
	"crdx.org/lighthouse/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	g := app.Group("/adapter/:id<int>").Use(auth.Admin)

	g.Post("/delete", Delete)
	g.Get("/edit", ViewEdit)
	g.Post("/edit", Edit)
}
