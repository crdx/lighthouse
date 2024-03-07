package admin

import (
	"crdx.org/lighthouse/controllers/admin/audit"
	"crdx.org/lighthouse/controllers/admin/mappings"
	"crdx.org/lighthouse/controllers/admin/settings"
	"crdx.org/lighthouse/controllers/admin/users"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/util"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin").
		Use(auth.Admin)

	adminGroup.Get("/", Index)

	adminGroup.Get("/settings", settings.List).Name("admin")
	adminGroup.Post("/settings", settings.Save)

	adminGroup.Get("/users", users.List).Name("admin")
	adminGroup.Get("/users/create", users.ViewCreate).Name("admin")
	adminGroup.Post("/users/create", users.Create)

	userGroup := adminGroup.Group("/users/:id<int>").
		Use(util.NewParseParam[m.User]("id", "user"))

	userGroup.Get("/edit", users.ViewEdit).Name("admin")
	userGroup.Post("/edit", users.Edit)
	userGroup.Post("/delete", users.Delete)
	userGroup.Post("/become", users.Become)

	auditGroup := adminGroup.Group("/audit")
	auditGroup.Get("/", audit.List).Name("admin")

	adminGroup.Get("/mappings", mappings.View).Name("admin")
	adminGroup.Post("/mappings", mappings.EditSources)

	adminGroup.Group("/mappings/:id<int>").
		Use(util.NewParseParam[m.Mapping]("id", "mapping")).
		Post("/delete", mappings.DeleteMapping)

	adminGroup.Post("/mappings/add", mappings.AddMapping)
	adminGroup.Get("/mappings/add", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/mappings")
	})
}
