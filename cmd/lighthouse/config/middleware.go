package config

import (
	"embed"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"crdx.org/lighthouse/pkg/middleware/minify"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func InitMiddleware(app *fiber.App, assets *embed.FS, dbConfig *db.Config) {
	app.Use(helmet.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(etag.New())

	app.Use(static.New("/assets", static.Config{
		FS: assets,
	}))

	if env.Production() {
		app.Use(recover.New(recover.Config{}))
	} else {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	app.Use(minify.New())

	app.Use(limiter.New(limiter.Config{
		Max:        300,
		Expiration: 60 * time.Second,
		// LimiterMiddleware: limiter.SlidingWindow{},
	}))

	if !env.Production() {
		app.Use(logger.New())
	}

	app.Use(NewSessionMiddleware(dbConfig))

	if env.DisableAuth() {
		app.Use(auth.AutoLogin(constants.RoleAdmin))
	} else {
		app.Use(auth.New())
	}
}
