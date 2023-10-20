package main

import (
	"crdx.org/lighthouse/controllers/activityController"
	"crdx.org/lighthouse/controllers/adapterController"
	"crdx.org/lighthouse/controllers/deviceController"
	"crdx.org/lighthouse/controllers/notificationController"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App) {
	deviceController.InitRoutes(app)
	adapterController.InitRoutes(app)
	activityController.InitRoutes(app)
	notificationController.InitRoutes(app)
}
