package conf

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
)

func GetDbConfig() *db.Config {
	return &db.Config{
		Name:          env.DatabaseName,
		User:          env.DatabaseUser,
		Pass:          env.DatabasePass,
		Host:          env.DatabaseHost,
		Socket:        env.DatabaseSocket,
		TimeZone:      env.DatabaseTimeZone,
		CharSet:       env.DatabaseCharSet,
		Models:        m.All,
		Migrations:    migrations,
		Colour:        !env.Production,
		Debug:         env.Debug,
		SlowThreshold: 250 * time.Millisecond,
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
		Models:     m.All,
		Migrations: migrations,
		Colour:     false,
		Debug:      false,
		Fresh:      true,
	}
}
