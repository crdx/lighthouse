package service

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/util"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	serviceGroup := app.Group("/service/:id<int>").
		Use(auth.Editor).
		Use(util.NewParseParam[m.Service]("id", "service"))

	serviceGroup.Post("/delete", Delete)
	serviceGroup.Get("/edit", ViewEdit).Name("devices")
	serviceGroup.Post("/edit", Edit)
}
