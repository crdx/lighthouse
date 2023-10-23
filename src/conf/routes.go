package conf

import (
	"crdx.org/lighthouse/controllers/activityController"
	"crdx.org/lighthouse/controllers/adapterController"
	"crdx.org/lighthouse/controllers/adminController"
	"crdx.org/lighthouse/controllers/deviceController"
	"crdx.org/lighthouse/controllers/notificationController"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	activityController.InitRoutes(app)
	adapterController.InitRoutes(app)
	adminController.InitRoutes(app)
	deviceController.InitRoutes(app)
	notificationController.InitRoutes(app)

	// Catch all requests not defined above.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}