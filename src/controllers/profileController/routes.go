package profileController

import (
	"crdx.org/lighthouse/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/profile", View)
	app.Post("/profile", auth.Editor, Edit)
}
