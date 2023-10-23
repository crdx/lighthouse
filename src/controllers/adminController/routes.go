package adminController

import (
	"crdx.org/lighthouse/controllers/adminController/settings"
	"crdx.org/lighthouse/controllers/adminController/users"
	"crdx.org/lighthouse/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	g := app.Group("/admin").Use(auth.Admin)

	g.Get("/", Index)

	g.Get("/settings", settings.List)
	g.Post("/settings", settings.Save)

	g.Get("/users", users.List)
	g.Get("/users/:id<int>/edit", users.ViewEdit)
	g.Post("/users/:id<int>/edit", users.Edit)
	g.Post("/users/:id<int>/delete", users.Delete)
}
