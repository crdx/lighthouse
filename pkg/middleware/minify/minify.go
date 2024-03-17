package minify

import (
	"crdx.org/lighthouse/pkg/util/webutil"
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := ctx.Next(); err != nil {
			return err
		}

		if !webutil.IsHTMLContentType(string(ctx.Response().Header.ContentType())) {
			return nil
		}

		if html, err := webutil.MinifyHTML(ctx.Response().Body()); err != nil {
			return err
		} else {
			ctx.Response().SetBody(html)
		}

		return nil
	}
}
