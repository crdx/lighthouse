package api

import (
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")

	apiGroup.Get("/icon/search", SearchIcon)
}
