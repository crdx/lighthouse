package main

import (
	"embed"
	"log"

	"crdx.org/db"
	"crdx.org/duckopt/v2"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/util"
	"crdx.org/lighthouse/util/mailutil"
	"crdx.org/lighthouse/util/timeutil"
	"crdx.org/session"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

//go:generate go run ../helpers/modelgen/main.go

//go:embed views/*
var views embed.FS

func getUsage() string {
	return `
        Usage:
            $0 [options] [--env PATH]

        Options:
            --env PATH    Read environment file
            -h, --help    Show help
    `
}

type Opts struct {
	EnvFile string `docopt:"--env"`
}

func main() {
	log.SetFlags(0)
	opts := duckopt.MustBind[Opts](getUsage(), "$0")

	initEnvironment(opts.EnvFile)

	dbConfig := conf.GetDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetSessionConfig(), dbConfig)

	app := fiber.New(conf.GetFiberConfig(views))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	if env.LiveReload() {
		app.Get("/hang", func(c *fiber.Ctx) error {
			select {}
		})
	}

	initMiddleware(app)
	conf.InitRoutes(app)

	initPackages()
	startServices()

	panic(app.Listen(env.Host() + ":" + env.Port()))
}

func initEnvironment(envFile string) {
	if envFile != "" {
		lo.Must0(env.InitFrom(envFile))
	} else if util.PathExists(".env") {
		lo.Must0(env.InitFrom(".env"))
	} else {
		env.Init()
	}

	if err := env.Validate(); err != nil {
		log.Fatalf("Error: failed to initialise environment:\n%s", err)
	}
}

func initPackages() {
	timeutil.Init(&timeutil.Config{
		Timezone: func() string { return settingR.Get(settingR.Timezone) },
	})

	mailutil.Init(&mailutil.Config{
		SendToStdErr: !env.Production(),
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
