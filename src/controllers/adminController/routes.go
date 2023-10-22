package adminController

import "github.com/gofiber/fiber/v2"

func InitRoutes(app *fiber.App) {
	app.Get("/admin", Index)
	app.Post("/admin/settings", SaveSettings)
}
