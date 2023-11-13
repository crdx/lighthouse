package main

import (
	"embed"
	"net/http"
	"time"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m/repo/userR"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/middleware/minify"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

//go:embed assets/*
var assets embed.FS

func initMiddleware(app *fiber.App) {
	app.Use(helmet.New())

	if env.Production() {
		app.Use(etag.New())
	}

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(assets),
		PathPrefix: "assets",
	}))

	if env.Production() {
		app.Use(recover.New(recover.Config{}))
	} else {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if env.Production() {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))

		app.Use(minify.New())

		app.Use(limiter.New(limiter.Config{
			Max:        300,
			Expiration: 60 * time.Second,
			// LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}

	if !env.Production() {
		app.Use(logger.New())
	}

	if env.DisableAuth() {
		app.Use(auth.AutoLogin(userR.RoleAdmin))
	} else {
		app.Use(auth.New())
	}
}
