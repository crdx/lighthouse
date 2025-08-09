package config

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers/seeder"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/schema"
	"crdx.org/lighthouse/pkg/env"
	"github.com/google/uuid"
)

func GetDbConfig() *db.Config {
	return &db.Config{
		Open: func(dsn *db.DSN) (*sql.DB, error) {
			return sql.Open("mysql", dsn.Format())
		},
		DataSource: db.NewDSN().Apply(func(dsn *db.DSN) *db.DSN {
			dsn.DBName = env.DatabaseName()
			dsn.Username = env.DatabaseUsername()
			dsn.Password = env.DatabasePassword()
			dsn.Address = env.DatabaseAddress()
			dsn.Protocol = env.DatabaseProtocol()
			return dsn
		}),
		Create:       true,
		Migrations:   schema.GetMigrations(),
		EnableLogger: env.Debug(),
		Seed:         seed,
	}
}

func GetTestDbConfig() *db.Config {
	return &db.Config{
		Open: func(dsn *db.DSN) (*sql.DB, error) {
			return sql.Open("mysql", dsn.Format())
		},
		DataSource: db.NewDSN().Apply(func(dsn *db.DSN) *db.DSN {
			dsn.DBName = env.DatabaseName() + "_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
			dsn.Username = env.DatabaseUsername()
			dsn.Password = env.DatabasePassword()
			dsn.Address = env.DatabaseAddress()
			dsn.Protocol = env.DatabaseProtocol()
			return dsn
		}),
		Migrations:   schema.GetMigrations(),
		Fresh:        true,
		EnableLogger: false,
		Seed:         seeder.Run,
	}
}
