package notifier_test

import (
	"log/slog"
	"net"
	"testing"
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier"
	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

// closedPort returns a local port that refuses connections.
func closedPort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	_, port, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	require.NoError(t, listener.Close())

	return port
}

func initMail(t *testing.T, config *mailutil.Config) {
	t.Helper()

	mailutil.Init(config)
	// Restore the configuration from helpers.TestMain.
	t.Cleanup(func() { mailutil.Init(&mailutil.Config{Enabled: func() bool { return false }}) })
}

func run(t *testing.T) {
	t.Helper()

	service := notifier.New()
	require.NoError(t, service.Init(&services.Args{Logger: slog.New(slog.DiscardHandler)}))
	require.NoError(t, service.Run())
}

// A failed send leaves the rows unprocessed. So the next run sends the notification again, and only
// a sent notification is recorded.
func TestFailedSendIsRetried(t *testing.T) {
	lo.Must(db.Exec("TRUNCATE device_discovery_notifications"))
	lo.Must(db.Exec("TRUNCATE device_limit_notifications"))
	lo.Must(db.Exec("TRUNCATE device_service_notifications"))
	lo.Must(db.Exec("TRUNCATE device_state_notifications"))
	lo.Must(db.Exec("TRUNCATE notifications"))

	device, found := db.FindDevice(4)
	require.True(t, found)
	db.CreateDeviceLimitNotification(&db.DeviceLimitNotification{
		DeviceID:       device.ID,
		StateUpdatedAt: db.Now().Add(-2 * time.Hour),
		Limit:          "1 hour",
	})

	port := closedPort(t)
	initMail(t, &mailutil.Config{
		Enabled:     func() bool { return true },
		Host:        func() string { return "127.0.0.1" },
		Port:        func() string { return port },
		User:        func() string { return "" },
		Pass:        func() string { return "" },
		FromAddress: func() string { return "from@example.test" },
		ToAddress:   func() string { return "to@example.test" },
		FromHeader:  func() string { return "From" },
		ToHeader:    func() string { return "To" },
	})
	run(t)

	assert.Len(t, db.FindUnprocessedDeviceLimitNotifications(), 1)
	assert.Empty(t, db.FindNotifications())

	initMail(t, &mailutil.Config{
		SendToStdErr: true,
		Enabled:      func() bool { return true },
		FromHeader:   func() string { return "From" },
		ToHeader:     func() string { return "To" },
		FromAddress:  func() string { return "from@example.test" },
		ToAddress:    func() string { return "to@example.test" },
	})
	run(t)

	assert.Empty(t, db.FindUnprocessedDeviceLimitNotifications())
	assert.Len(t, db.FindNotifications(), 1)
}
