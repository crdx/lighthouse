package activity

import (
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	app.Get("/activity", List).Name("activity")
}
