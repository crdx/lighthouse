package helpers

import (
	"os"
	"testing"

	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"crdx.org/session"
	"github.com/samber/lo"
)

// Start begins a new db transaction and returns a func that will roll back the transaction. This
// function is NOT thread-safe, but this is fine because each test within a test package runs in
// serial.
func Start() func() {
	instance := db.Instance()
	db.SetInstance(instance.Begin())

	return func() {
		db.Instance().Rollback()
		db.SetInstance(instance)
	}
}

// TestMain initialises the environment and database for the current test package, runs the tests,
// and then drops the test database.
func TestMain(m *testing.M) {
	env.Init()

	dbConfig := conf.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetTestSessionConfig(), dbConfig)

	timeutil.Init(&timeutil.Config{Timezone: func() string { return "Europe/London" }})
	mailutil.Init(&mailutil.Config{Enabled: func() bool { return false }})

	exitCode := m.Run()

	db.Exec("DROP DATABASE " + dbConfig.Name)
	os.Exit(exitCode)
}
