package minify

import (
	"crdx.org/lighthouse/pkg/util/webutil"
	"github.com/gofiber/fiber/v3"
)

func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}

		if !webutil.IsHTMLContentType(string(c.Response().Header.ContentType())) {
			return nil
		}

		if html, err := webutil.MinifyHTML(c.Response().Body()); err != nil {
			return err
		} else {
			c.Response().SetBody(html)
		}

		return nil
	}
}
