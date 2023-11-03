package apiController

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")

	apiGroup.Get("icon/search", SearchIcon)
}
