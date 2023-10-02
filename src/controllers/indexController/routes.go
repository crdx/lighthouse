package indexController

import "github.com/gofiber/fiber/v2"

func InitRoutes(app *fiber.App) {
	app.Get("/", Get)
	app.Get("/device/:id<int>", ViewDevice)
}
