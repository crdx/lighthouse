package indexController

import (
	"crdx.org/lighthouse/models/deviceModel"

	"github.com/gofiber/fiber/v2"
)

// var log = logger.New().With("controller", "index")

func Get(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"devices": deviceModel.GetListView(),
	})
}
