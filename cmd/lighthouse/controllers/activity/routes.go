package activity

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/activity", List).Name("activity")
}
