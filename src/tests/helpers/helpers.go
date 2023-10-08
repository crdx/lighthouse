package helpers

import (
	"encoding/gob"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func Init() {
	dbConfig := conf.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetTestSessionConfig(), dbConfig)
	gob.Register(&flash.Message{})

	Seed()
}

func App() *fiber.App {
	return fiber.New(conf.GetTestFiberConfig())
}
