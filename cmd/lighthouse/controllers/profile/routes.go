package profile

import (
	"crdx.org/lighthouse/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/profile", View).Name("profile")
	app.Post("/profile", auth.Editor, Edit)
}
