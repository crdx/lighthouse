package conf

import (
	"crdx.org/lighthouse/controllers/activity"
	"crdx.org/lighthouse/controllers/adapter"
	"crdx.org/lighthouse/controllers/admin"
	"crdx.org/lighthouse/controllers/api"
	"crdx.org/lighthouse/controllers/device"
	"crdx.org/lighthouse/controllers/notification"
	"crdx.org/lighthouse/controllers/profile"
	"crdx.org/lighthouse/controllers/service"
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
