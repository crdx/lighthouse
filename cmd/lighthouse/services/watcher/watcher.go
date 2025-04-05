// Package watcher watches for devices whose state has changed.
package watcher

import (
	"log/slog"
	"time"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceR"
)

type Watcher struct {
	logger                 *slog.Logger
	limitNotificationCache map[int64]time.Time
}

func New() *Watcher {
	return &Watcher{}
}

func (self *Watcher) Init(args *services.Args) error {
	self.logger = args.Logger
	self.limitNotificationCache = map[int64]time.Time{}

	return nil
}

func (self *Watcher) Run() error {
	for _, device := range db.FindDevicesSorted() {
		if device.Origin {
			continue
		}

		gracePeriod := device.GracePeriodDuration()

		logger := self.logger.With(slog.Group(
			"device",
			"id", device.ID,
			"name", device.Name,
		))

		threshold := db.Now().Add(-gracePeriod)

		if device.LastSeenAt.Before(threshold) {
			if device.State == deviceR.StateOnline {
				deviceOffline(device, threshold)
				logger.Info("device is offline")
			}
		} else {
			if device.State == deviceR.StateOffline {
				deviceOnline(device)
				logger.Info("device is online")
			}
		}

		limit := device.LimitDuration()

		if limit == 0 {
			continue
		}

		if device.State == deviceR.StateOnline && device.StateUpdatedAt.Before(db.Now().Add(-limit)) {
			if self.limitNotificationCache[device.ID].Equal(device.StateUpdatedAt) {
				continue
			}

			self.limitNotificationCache[device.ID] = device.StateUpdatedAt

			if db.CountPreviousDeviceLimitNotifications(device.ID) > 0 {
				continue
			}

			db.CreateDeviceLimitNotification(&db.DeviceLimitNotification{
				DeviceID:       device.ID,
				StateUpdatedAt: device.StateUpdatedAt,
				Limit:          device.Limit,
			})

			logger.Info("device overstayed its welcome", "limit", limit)
		}
	}

	return nil
}

func deviceOnline(device *db.Device) {
	state := deviceR.StateOnline

	db.CreateDeviceStateLog(&db.DeviceStateLog{
		DeviceID:    device.ID,
		State:       state,
		GracePeriod: device.GracePeriod,
	})

	if device.Watch {
		db.CreateDeviceStateNotification(&db.DeviceStateNotification{
			DeviceID:    device.ID,
			State:       state,
			GracePeriod: device.GracePeriod,
		})
	}

	device.UpdateState(state)
	device.UpdateStateUpdatedAt(db.Now())
}

func deviceOffline(device *db.Device, transitionedAt time.Time) {
	state := deviceR.StateOffline

	db.CreateDeviceStateLog(&db.DeviceStateLog{
		DeviceID:    device.ID,
		State:       state,
		GracePeriod: device.GracePeriod,
		CreatedAt:   transitionedAt,
		UpdatedAt:   db.N(transitionedAt),
	})

	if device.Watch {
		db.CreateDeviceStateNotification(&db.DeviceStateNotification{
			DeviceID:    device.ID,
			State:       state,
			GracePeriod: device.GracePeriod,
			CreatedAt:   transitionedAt,
			UpdatedAt:   db.N(transitionedAt),
		})
	}

	device.UpdateState(state)
	device.UpdateStateUpdatedAt(db.Now())
}
