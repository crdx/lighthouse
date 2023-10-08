package adapterController

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	adapterGroup := app.Group("/adapter/:id<int>")
	adapterGroup.Post("/delete", Delete)
	adapterGroup.Get("/edit", ViewEdit)
	adapterGroup.Post("/edit", Edit)
}
