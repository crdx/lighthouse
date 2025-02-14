package adapter

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	adapterGroup := app.Group("/adapter/:id<int>").
		Use(auth.Editor).
		Use(parseparam.New("id", db.FindAdapter))

	adapterGroup.Post("/delete", Delete)
	adapterGroup.Get("/edit", ViewEdit).Name("devices")
	adapterGroup.Post("/edit", Edit)
}
