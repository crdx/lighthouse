package adminController

import (
	"crdx.org/lighthouse/controllers/adminController/settings"
	"crdx.org/lighthouse/controllers/adminController/users"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/admin", Index)

	app.Get("/admin/settings", settings.List)
	app.Post("/admin/settings", settings.Save)

	app.Get("/admin/users", users.List)
	app.Get("/admin/users/:id<int>/edit", users.ViewEdit)
	app.Post("/admin/users/:id<int>/edit", users.Edit)
	app.Post("/admin/users/:id<int>/delete", users.Delete)
}
