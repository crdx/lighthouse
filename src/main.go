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

	initPackages()
	startServices()

	panic(app.Listen(env.BindHost + ":" + env.BindPort))
}

func initPackages() {
	timeutil.Init(&timeutil.Config{
		Timezone: func() string { return settingR.Get(settingR.Timezone) },
	})

	mailutil.Init(&mailutil.Config{
		SendToStdErr: !env.Production,
		Enabled:      func() bool { return settingR.GetBool(settingR.EnableMail) },
		Host:         func() string { return settingR.Get(settingR.SMTPHost) },
		Port:         func() string { return settingR.Get(settingR.SMTPPort) },
		User:         func() string { return settingR.Get(settingR.SMTPUser) },
		Pass:         func() string { return settingR.Get(settingR.SMTPPass) },
		FromAddress:  func() string { return settingR.Get(settingR.MailFromAddress) },
		ToAddress:    func() string { return settingR.Get(settingR.MailToAddress) },
		FromHeader:   func() string { return settingR.Get(settingR.MailFromHeader) },
		ToHeader:     func() string { return settingR.Get(settingR.MailToHeader) },
	})
}
