package main

import (
	"embed"
	"log"
	"time"

	"crdx.org/duckopt/v2"
	"crdx.org/lighthouse/cmd/lighthouse/config"
	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier"
	"crdx.org/lighthouse/cmd/lighthouse/services/pinger"
	"crdx.org/lighthouse/cmd/lighthouse/services/prober"
	"crdx.org/lighthouse/cmd/lighthouse/services/reader"
	"crdx.org/lighthouse/cmd/lighthouse/services/vendor"
	"crdx.org/lighthouse/cmd/lighthouse/services/watcher"
	"crdx.org/lighthouse/cmd/lighthouse/services/writer"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/settingR"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/logger"
	"crdx.org/lighthouse/pkg/util"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"crdx.org/session/v3"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

//go:embed assets/*
var assets embed.FS

//go:embed views/*
var views embed.FS

func getUsage() string {
	return `
		Usage:
			$0 [options] [--env PATH]

		Options:
			--env PATH    Read environment file
	`
}

type Opts struct {
	EnvFile string `docopt:"--env"`
}

func main() {
	log.SetFlags(0)
	opts := duckopt.MustBind[Opts](getUsage(), "$0")

	initEnvironment(opts.EnvFile)

	dbConfig := config.GetDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(config.GetSessionConfig(), dbConfig.DataSource.Format())

	app := fiber.New(config.GetFiberConfig(views))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	if env.LiveReload() {
		app.Get("/hang", func(c fiber.Ctx) error {
			select {}
		})
	}

	logger.Init()

	config.InitMiddleware(app, &assets)
	config.InitRoutes(app)

	initPackages()

	if !env.DisableServices() {
		startServices()
	}

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
		log.Fatal(err)
	}
}

func initPackages() {
	timeutil.Init(&timeutil.Config{
		Timezone: settingR.Timezone,
	})

	mailutil.Init(&mailutil.Config{
		SendToStdErr: !env.Production(),
		Enabled:      settingR.EnableMail,
		Host:         settingR.SMTPHost,
		Port:         settingR.SMTPPort,
		User:         settingR.SMTPUser,
		Pass:         settingR.SMTPPass,
		FromAddress:  settingR.MailFromAddress,
		ToAddress:    settingR.MailToAddress,
		FromHeader:   settingR.MailFromHeader,
		ToHeader:     settingR.MailToHeader,
	})
}

func startServices() {
	services.Start("writer", &services.Config{
		Service:         writer.New(),
		RunIntervalFunc: settingR.DeviceScanInterval,
		Enabled:         settingR.EnableDeviceScan,
	})

	services.Start("reader", &services.Config{
		Service: reader.New(), // Reader service runs forever.
	})

	services.Start("vendor", &services.Config{
		Service:     vendor.New(),
		RunInterval: 5 * time.Second,
		Enabled:     func() bool { return settingR.MACVendorsAPIKey() != "" },
		Quiet:       true,
	})

	watcherStartDelay := 30 * time.Second

	services.Start("watcher", &services.Config{
		Service:     watcher.New(),
		RunInterval: 10 * time.Second,
		StartDelay:  watcherStartDelay,
		Quiet:       true,
	})

	services.Start("pinger", &services.Config{
		Service:         pinger.New(),
		RunIntervalFunc: settingR.DeviceScanInterval,
		Enabled:         settingR.EnableDeviceScan,

		// Give it time to establish the current state of devices.
		StartDelay: watcherStartDelay + (30 * time.Second),
	})

	services.Start("notifier", &services.Config{
		Service:     notifier.New(),
		RunInterval: 1 * time.Minute,

		// Give it time to establish the current state of devices.
		StartDelay: watcherStartDelay + (30 * time.Second),

		// A fairly high initial restart interval as we don't want to risk spamming notifications.
		InitialRestartInterval: 1 * time.Minute,
	})

	services.Start("prober", &services.Config{
		Service:    prober.New(),
		StartDelay: 5 * time.Minute,

		RunIntervalFunc: settingR.ServiceScanInterval,
		Enabled:         settingR.EnableServiceScan,
	})
}
