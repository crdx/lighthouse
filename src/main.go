package main

import (
	"embed"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/session"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	if env.EnableLiveReload {
		app.Get("/hang", func(c *fiber.Ctx) error {
			select {}
		})
	}

	initMiddleware(app)
	initRoutes(app)

	// Catch all requests not defined in initRoutes above.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	startServices()
	registerValidators()

	panic(app.Listen(env.BindHost + ":" + env.BindPort))
}

func registerValidators() {
	validate.Register("timezone", "invalid timezone", func(value string) bool {
		_, err := time.LoadLocation(value)
		return err == nil
	})
}
