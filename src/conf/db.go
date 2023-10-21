package conf

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/setting"
	"crdx.org/lighthouse/tests/helpers/seeder"
)

func seed() error {
	db.Save(&m.Setting{Name: setting.Timezone, Value: "Europe/London"})
	db.Save(&m.Setting{Name: setting.Watch, Value: "1"})
	return nil
}

func GetDbConfig() *db.Config {
	return &db.Config{
		Name:          env.DatabaseName,
		User:          env.DatabaseUser,
		Pass:          env.DatabasePass,
		Host:          env.DatabaseHost,
		Socket:        env.DatabaseSocket,
		TimeZone:      env.DatabaseTimeZone,
		CharSet:       env.DatabaseCharSet,
		Models:        models,
		Migrations:    migrations,
		Colour:        !env.Production,
		Debug:         env.Debug,
		SlowThreshold: 250 * time.Millisecond,
		Seed:          seed,
	}
}

func GetTestDbConfig() *db.Config {
	return &db.Config{
		Name:       env.DatabaseName + "_test",
		User:       env.DatabaseUser,
		Pass:       env.DatabasePass,
		Host:       env.DatabaseHost,
		Socket:     env.DatabaseSocket,
		TimeZone:   env.DatabaseTimeZone,
		CharSet:    env.DatabaseCharSet,
		Models:     models,
		Migrations: migrations,
		Colour:     false,
		Debug:      false,
		Fresh:      true,
		Seed:       seeder.Run,
	}
}
