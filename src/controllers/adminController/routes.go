package adminController

import (
	"crdx.org/lighthouse/controllers/adminController/audit"
	"crdx.org/lighthouse/controllers/adminController/settings"
	"crdx.org/lighthouse/controllers/adminController/users"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/util"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin").
		Use(auth.Admin)

	adminGroup.Get("/", Index)

	adminGroup.Get("/settings", settings.List)
	adminGroup.Post("/settings", settings.Save)

	adminGroup.Get("/users", users.List)
	adminGroup.Get("/users/create", users.ViewCreate)
	adminGroup.Post("/users/create", users.Create)

	userGroup := adminGroup.Group("/users/:id<int>").
		Use(util.NewParseParam[m.User]("id", "user"))

	userGroup.Get("/edit", users.ViewEdit)
	userGroup.Post("/edit", users.Edit)
	userGroup.Post("/delete", users.Delete)

	auditGroup := adminGroup.Group("/audit")
	auditGroup.Get("/", audit.List)
}
