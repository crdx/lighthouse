package api

import (
	"crdx.org/lighthouse/pkg/fontawesome"
	"github.com/gofiber/fiber/v3"
)

func SearchIcon(c fiber.Ctx) error {
	icons, hasMore, err := fontawesome.Search(c.Query("q"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.JSON(fiber.Map{
		"icons":   icons,
		"hasMore": hasMore,
	})
}
