package notification

import (
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	app.Get("/notifications", List).Name("notifications")
}
