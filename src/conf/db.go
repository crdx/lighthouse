package conf

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/tests/helpers/seeder"
	"crdx.org/lighthouse/util/stringutil"
)

func seed() error {
	db.Save(&m.Setting{Name: "timezone", Value: "Europe/London"})
	db.Save(&m.Setting{Name: "watch", Value: "1"})

	db.Save(&m.User{Username: "root", PasswordHash: stringutil.Hash(env.DefaultRootPass()), Admin: true})
	db.Save(&m.User{Username: "anon", PasswordHash: stringutil.Hash(env.DefaultAnonPass()), Admin: false})

	return nil
}

func GetDbConfig() *db.Config {
	return &db.Config{
		Name:          env.DatabaseName(),
		User:          env.DatabaseUser(),
		Pass:          env.DatabasePass(),
		Host:          env.DatabaseHost(),
		Socket:        env.DatabaseSocket(),
		TimeZone:      env.DatabaseTimezone(),
		CharSet:       env.DatabaseCharset(),
		Models:        models,
		Migrations:    migrations,
		Colour:        !env.Production(),
		Debug:         env.Debug(),
		SlowThreshold: 250 * time.Millisecond,
		Seed:          seed,
	}
}

func GetTestDbConfig() *db.Config {
	return &db.Config{
		Name:       env.DatabaseName() + "_test",
		User:       env.DatabaseUser(),
		Pass:       env.DatabasePass(),
		Host:       env.DatabaseHost(),
		Socket:     env.DatabaseSocket(),
		TimeZone:   env.DatabaseTimezone(),
		CharSet:    env.DatabaseCharset(),
		Models:     models,
		Migrations: migrations,
		Colour:     false,
		Debug:      false,
		Fresh:      true,
		Seed:       seeder.Run,
	}
}
