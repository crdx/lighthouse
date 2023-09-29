package main

import (
	"embed"
	"net/http"
	"time"

	"crdx.org/lighthouse/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

//go:embed assets/*
var assets embed.FS

func initMiddleware(app *fiber.App) {
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(assets),
		PathPrefix: "assets",
	}))

	app.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(assets),
		File:       "assets/favicon.svg",
	}))

	if env.AuthType == env.AuthTypeBasic {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				env.AuthUser: env.AuthPass,
			},
		}))
	}

	if env.Production {
		app.Use(recover.New(recover.Config{}))
	} else {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if env.Production {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))

		app.Use(limiter.New(limiter.Config{
			Max:        300,
			Expiration: 60 * time.Second,
			// LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}

	if !env.Production {
		app.Use(logger.New())
	}
}
