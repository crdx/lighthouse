package main

import (
	"embed"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
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

	initMiddleware(app)
	initRoutes(app)

	// Catch all requests not defined in initRoutes above.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	startServices()

	panic(app.Listen(env.BindHost + ":" + env.BindPort))
}
