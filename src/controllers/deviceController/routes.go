package deviceController

import (
	"crdx.org/lighthouse/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/", List)

	g := app.Group("/device/:id<int>")

	g.Get("/", View)

	g.Post("/delete", auth.Admin, Delete)
	g.Get("/edit", auth.Admin, ViewEdit)
	g.Post("/edit", auth.Admin, Edit)
	g.Post("/merge", auth.Admin, Merge)
}
