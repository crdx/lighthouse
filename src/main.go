package main

import (
	"bytes"
	"embed"
	"encoding/gob"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/util/webutil"
	"crdx.org/session"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

//go:generate go run ../helpers/modelgen/main.go

//go:embed views/*
var views embed.FS

func main() {
	env.Check()

	dbConfig := conf.GetDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetSessionConfig(), dbConfig)

	app := fiber.New(conf.GetFiberConfig(views))

	if env.Production {
		initMinifier(app)
	}

	initHealthCheck(app)
	initMiddleware(app)
	initFlash(app)
	initRoutes(app)

	// Catch all requests not defined in initRoutes above.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	startServices()

	panic(app.Listen(env.BindHost + ":" + env.BindPort))
}

func initMinifier(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err != nil {
			return err
		}

		if !webutil.IsHTMLContentType(string(ctx.Response().Header.ContentType())) {
			return nil
		}

		htmlMinifier := &html.Minifier{}
		htmlMinifier.KeepComments = false            // Preserve all comments
		htmlMinifier.KeepConditionalComments = false // Preserve all IE conditional comments
		htmlMinifier.KeepDefaultAttrVals = false     // Preserve default attribute values
		htmlMinifier.KeepDocumentTags = false        // Preserve html, head and body tags
		htmlMinifier.KeepEndTags = false             // Preserve all end tags
		htmlMinifier.KeepWhitespace = false          // Preserve whitespace characters but still collapse multiple into one
		htmlMinifier.KeepQuotes = false              // Preserve quotes around attribute values

		var minifiedBody bytes.Buffer
		htmlMinifier.Minify(
			minify.New(),
			&minifiedBody,
			bytes.NewReader(ctx.Response().Body()),
			nil,
		)

		ctx.Response().SetBody(minifiedBody.Bytes())

		return nil
	})
}

func initHealthCheck(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}

func initFlash(app *fiber.App) {
	gob.Register(flash.Message{})

	app.Use(func(c *fiber.Ctx) error {
		if flashMessage, found := session.GetOnce[flash.Message](c, flash.Key); found {
			c.Locals(flash.Key, flashMessage)
		}

		return c.Next()
	})
}
