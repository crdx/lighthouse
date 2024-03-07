package profile

import (
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

func View(c *fiber.Ctx) error {
	return c.Render("profile/view", fiber.Map{
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}
