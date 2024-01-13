package conf

import (
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/tests/helpers/seeder"
	"github.com/google/uuid"
)

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
		Migrations:    dbMigrations,
		Colour:        !env.Production(),
		Debug:         env.Debug(),
		SlowThreshold: 250 * time.Millisecond,
		Seed:          seed,
	}
}

func GetTestDbConfig() *db.Config {
	return &db.Config{
		Name:       env.DatabaseName() + "_test_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
		User:       env.DatabaseUser(),
		Pass:       env.DatabasePass(),
		Host:       env.DatabaseHost(),
		Socket:     env.DatabaseSocket(),
		TimeZone:   env.DatabaseTimezone(),
		CharSet:    env.DatabaseCharset(),
		Models:     models,
		Migrations: dbMigrations,
		Colour:     false,
		Debug:      false,
		Fresh:      true,
		Seed:       seeder.Run,
	}
}
