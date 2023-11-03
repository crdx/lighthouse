package apiController

import (
	"crdx.org/lighthouse/pkg/fontawesome"
	"github.com/gofiber/fiber/v2"
)

func SearchIcon(c *fiber.Ctx) error {
	icons, hasMore := fontawesome.Search(c.Query("q"))

	return c.JSON(fiber.Map{
		"icons":   icons,
		"hasMore": hasMore,
	})
}
