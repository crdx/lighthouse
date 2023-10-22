package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func Init(state auth.State) *Session {
	dbConfig := conf.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetTestSessionConfig(), dbConfig)


	app := fiber.New(conf.GetTestFiberConfig())

	if state == auth.StateUnauthenticated {
		app.Use(auth.New())
	} else {
		app.Use(auth.AutoLogin(state))
	}

	conf.InitRoutes(app)
	return NewSession(app)
}
