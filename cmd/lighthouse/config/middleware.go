package config

import (
	"io/fs"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func InitMiddleware(app *fiber.App, assets fs.FS, dbConfig *db.Config) {
	if !env.Production() {
		app.Use(logger.New())
	}

	app.Use(limiter.New(limiter.Config{Max: 300, Expiration: 60 * time.Second}))
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New())

	app.Use(static.New("/assets", static.Config{FS: assets}))
	app.Use(NewSession(dbConfig.DataSource.Format()))

	if env.DisableAuth() {
		app.Use(auth.AutoLogin(constants.RoleAdmin))
	} else {
		app.Use(auth.New())
	}
}
