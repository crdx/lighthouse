package main

import (
	"embed"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/util/mailutil"
	"crdx.org/lighthouse/util/timeutil"
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
	conf.InitRoutes(app)

	InitMail()
	InitTime()

	startServices()

	panic(app.Listen(env.BindHost + ":" + env.BindPort))
}

func InitTime() {
	timeutil.Init(&timeutil.Config{
		Timezone: settingR.Get(settingR.Timezone),
	})
}

func InitMail() {
	mailutil.Init(&mailutil.Config{
		Enable:      settingR.GetBool(settingR.EnableMail),
		Host:        settingR.Get(settingR.SMTPHost),
		Port:        settingR.Get(settingR.SMTPPort),
		User:        settingR.Get(settingR.SMTPUser),
		Pass:        settingR.Get(settingR.SMTPPass),
		FromAddress: settingR.Get(settingR.MailFromAddress),
		ToAddress:   settingR.Get(settingR.MailToAddress),
		FromHeader:  settingR.Get(settingR.MailFromHeader),
		ToHeader:    settingR.Get(settingR.MailToHeader),
	})
}
