package config

import (
	"crdx.org/lighthouse/cmd/lighthouse/controllers/activity"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/adapter"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/admin"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/api"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/device"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/notification"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/profile"
	"crdx.org/lighthouse/cmd/lighthouse/controllers/service"
	"crdx.org/lighthouse/pkg/env"
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	if env.LiveReload() {
		app.Get("/hang", func(c fiber.Ctx) error {
			select {}
		})
	}

	activity.InitRoutes(app)
	adapter.InitRoutes(app)
	admin.InitRoutes(app)
	api.InitRoutes(app)
	device.InitRoutes(app)
	notification.InitRoutes(app)
	profile.InitRoutes(app)
	service.InitRoutes(app)
}
