package main

import (
	"crdx.org/lighthouse/controllers/adapterController"
	"crdx.org/lighthouse/controllers/deviceController"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App) {
	deviceController.InitRoutes(app)
	adapterController.InitRoutes(app)
}
