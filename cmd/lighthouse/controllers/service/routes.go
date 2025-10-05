package service

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	serviceGroup := app.Group("/service/:id<int>").
		Use(auth.Editor).
		Use(parseparam.New("id", db.FindService))

	serviceGroup.Post("/delete", Delete)
	serviceGroup.Get("/edit", ViewEdit).Name("devices")
	serviceGroup.Post("/edit", Edit)
}
