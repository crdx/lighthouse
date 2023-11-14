package profileController

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/profile", View)
	app.Post("/profile", Edit)
}
