package conf

import (
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"crdx.org/lighthouse/tests/helpers/seeder"
	"github.com/google/uuid"
)

func seed() error {
	db.Save(&m.Setting{Name: "timezone", Value: "Europe/London"})
	db.Save(&m.Setting{Name: "watch", Value: "1"})
	db.Save(&m.Setting{Name: "scan_interval", Value: "1 min"})

	db.Save(&m.User{Username: "root", PasswordHash: stringutil.Hash(env.DefaultRootPass()), Role: constants.RoleAdmin})
	db.Save(&m.User{Username: "anon", PasswordHash: stringutil.Hash(env.DefaultAnonPass()), Role: constants.RoleViewer})

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
