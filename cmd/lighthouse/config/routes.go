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
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	activity.InitRoutes(app)
	adapter.InitRoutes(app)
	admin.InitRoutes(app)
	api.InitRoutes(app)
	device.InitRoutes(app)
	notification.InitRoutes(app)
	profile.InitRoutes(app)
	service.InitRoutes(app)

	// Catch all requests not defined above.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
