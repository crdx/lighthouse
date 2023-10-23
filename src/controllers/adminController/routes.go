package adminController

import (
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/admin", Index)

	app.Get("/admin/settings", ListSettings)
	app.Post("/admin/settings", SaveSettings)

	app.Get("/admin/users", ListUsers)
}
