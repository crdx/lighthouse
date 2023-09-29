package main

import (
	"crdx.org/lighthouse/controllers/indexController"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App) {
	indexController.InitRoutes(app)
}
