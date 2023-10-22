package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func Init() {
	dbConfig := conf.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetTestSessionConfig(), dbConfig)
}

func App(state auth.State) *fiber.App {
	app := fiber.New(conf.GetTestFiberConfig())

	if state == auth.StateUnauthenticated {
		app.Use(auth.New())
	} else {
		app.Use(auth.AutoLogin(state))
	}

	return app
}
