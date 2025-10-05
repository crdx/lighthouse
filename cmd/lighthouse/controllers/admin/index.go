package admin

import (
	"github.com/gofiber/fiber/v3"
)

func Index(c fiber.Ctx) error {
	return c.Redirect().To("/admin/users")
}
