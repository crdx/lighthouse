package helpers

import (
	"os"
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/config"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"github.com/samber/lo"
)

var dbConfig *db.Config

// Start begins a new db transaction and registers a cleanup function that will roll back the
// transaction. This function is NOT thread-safe, but this is fine because each test within a test
// package runs in serial.
func Start(tb testing.TB) {
	tb.Helper()

	_ = db.BeginTransaction()

	tb.Cleanup(func() {
		_ = db.RollbackTransaction()
	})
}

// TestMain initialises the environment and database for the current test package, runs the tests,
// and then drops the test database.
func TestMain(m *testing.M) {
	env.Init()

	dbConfig = config.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))

	timeutil.Init(&timeutil.Config{Timezone: func() string { return "Europe/London" }})
	mailutil.Init(&mailutil.Config{Enabled: func() bool { return false }})

	exitCode := m.Run()

	_, _ = db.Exec("DROP DATABASE " + dbConfig.DataSource.DBName)
	os.Exit(exitCode)
}
