package watcher_test

import (
	"log/slog"
	"testing"
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/cmd/lighthouse/services/watcher"
	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceR"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestState(t *testing.T) {
	lo.Must(db.Exec("TRUNCATE device_limit_notifications"))
	lo.Must(db.Exec("TRUNCATE device_state_logs"))
	lo.Must(db.Exec("TRUNCATE device_state_notifications"))

	device, _ := db.FindDevice(1)
	device.UpdateWatch(true)

	device, _ = db.FindDevice(4)
	device.UpdateWatch(true)

	w := watcher.New()
	require.NoError(t, w.Init(&services.Args{Logger: slog.New(slog.DiscardHandler)}))
	require.NoError(t, w.Run())

	assert.Len(t, db.FindDeviceStateLogs(), 3)

	var zero int64
	assert.Greater(t, db.CountDeviceStateLogsForTest(1, deviceR.StateOffline), zero)
	assert.Greater(t, db.CountDeviceStateLogsForTest(2, deviceR.StateOffline), zero)
	assert.Greater(t, db.CountDeviceStateLogsForTest(4, deviceR.StateOnline), zero)

	assert.Greater(t, db.CountDeviceStateNotificationsForTest(1, deviceR.StateOffline), zero)
	assert.Greater(t, db.CountDeviceStateNotificationsForTest(4, deviceR.StateOnline), zero)
}

func TestLimit(t *testing.T) {
	lo.Must(db.Exec("TRUNCATE device_limit_notifications"))

	device, _ := db.FindDevice(4)
	device.UpdateLimit("1 hour")
	device.UpdateState(deviceR.StateOnline)
	device.UpdateStateUpdatedAt(db.Now().Add(-1 * time.Hour).Add(-5 * time.Minute))

	w := watcher.New()
	require.NoError(t, w.Init(&services.Args{Logger: slog.New(slog.DiscardHandler)}))
	require.NoError(t, w.Run())

	assert.Len(t, db.FindDeviceLimitNotifications(), 1)
}
