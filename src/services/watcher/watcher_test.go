package watcher_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/watcher"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestState(t *testing.T) {
	db.Exec("TRUNCATE device_limit_notifications")
	db.Exec("TRUNCATE device_state_logs")
	db.Exec("TRUNCATE device_state_notifications")

	db.For[m.Device](1).Update("watch", true)
	db.For[m.Device](4).Update("watch", true)

	w := watcher.New()
	require.NoError(t, w.Init(&services.Args{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}))
	require.NoError(t, w.Run())

	assert.Equal(t, int64(3), db.B[m.DeviceStateLog]().Count())
	assert.True(t, db.B[m.DeviceStateLog]("device_id = 1 and state = ?", deviceR.StateOffline).Exists())
	assert.True(t, db.B[m.DeviceStateLog]("device_id = 2 and state = ?", deviceR.StateOffline).Exists())
	assert.True(t, db.B[m.DeviceStateLog]("device_id = 4 and state = ?", deviceR.StateOnline).Exists())

	assert.True(t, db.B[m.DeviceStateNotification]("device_id = 1 and state = ? and processed = 0", deviceR.StateOffline).Exists())
	assert.True(t, db.B[m.DeviceStateNotification]("device_id = 4 and state = ? and processed = 0", deviceR.StateOnline).Exists())
}

func TestLimit(t *testing.T) {
	db.Exec("TRUNCATE device_limit_notifications")

	db.For[m.Device](1).Update(
		"limit", "1 hour",
		"state", deviceR.StateOnline,
		"state_updated_at", time.Now().Add(-1*time.Hour).Add(-5*time.Minute),
	)

	w := watcher.New()
	require.NoError(t, w.Init(&services.Args{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}))
	require.NoError(t, w.Run())

	assert.Equal(t, int64(1), db.B[m.DeviceLimitNotification]().Count())
}
