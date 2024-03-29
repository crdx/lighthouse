package helpers

import (
	"os"
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/config"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"crdx.org/session/v2"
	"github.com/samber/lo"
)

// Start begins a new db transaction and returns a func that will roll back the transaction. This
// function is NOT thread-safe, but this is fine because each test within a test package runs in
// serial.
func Start() func() {
	_ = db.BeginTransaction()

	return func() {
		_ = db.RollbackTransaction()
	}
}

// TestMain initialises the environment and database for the current test package, runs the tests,
// and then drops the test database.
func TestMain(m *testing.M) {
	env.Init()

	dbConfig := config.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(config.GetTestSessionConfig(), dbConfig.DataSource.Format())

	timeutil.Init(&timeutil.Config{Timezone: func() string { return "Europe/London" }})
	mailutil.Init(&mailutil.Config{Enabled: func() bool { return false }})

	exitCode := m.Run()

	_, _ = db.Exec("DROP DATABASE " + dbConfig.DataSource.DBName)
	os.Exit(exitCode)
}
