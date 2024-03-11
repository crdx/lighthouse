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
	config := db.Config{
		Name:          env.DatabaseName(),
		User:          env.DatabaseUsername(),
		Pass:          env.DatabasePassword(),
		TimeZone:      env.DatabaseTimezone(),
		CharSet:       env.DatabaseCharset(),
		Models:        models,
		Migrations:    dbMigrations,
		Colour:        !env.Production(),
		Debug:         env.Debug(),
		SlowThreshold: 250 * time.Millisecond,
		Seed:          seed,
	}

	if env.DatabaseProtocol() == "" || env.DatabaseProtocol() == "tcp" {
		config.Host = env.DatabaseAddress()
	} else if env.DatabaseProtocol() == "unix" {
		config.Socket = env.DatabaseAddress()
	}

	return &config
}

func GetTestDbConfig() *db.Config {
	config := db.Config{
		Name:       env.DatabaseName() + "_test_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
		User:       env.DatabaseUsername(),
		Pass:       env.DatabasePassword(),
		TimeZone:   env.DatabaseTimezone(),
		CharSet:    env.DatabaseCharset(),
		Models:     models,
		Migrations: dbMigrations,
		Colour:     false,
		Debug:      false,
		Fresh:      true,
		Seed:       seeder.Run,
	}

	if env.DatabaseProtocol() == "" || env.DatabaseProtocol() == "tcp" {
		config.Host = env.DatabaseAddress()
	} else if env.DatabaseProtocol() == "unix" {
		config.Socket = env.DatabaseAddress()
	}

	return &config
}
