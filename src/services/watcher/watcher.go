// Package watcher watches for devices whose state has changed.
package watcher

import (
	"log/slog"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/services"
)

type Watcher struct {
	logger                 *slog.Logger
	limitNotificationCache map[uint]time.Time
}

func New() *Watcher {
	return &Watcher{}
}

func (self *Watcher) Init(args *services.Args) error {
	self.logger = args.Logger
	self.limitNotificationCache = map[uint]time.Time{}

	return nil
}

func (self *Watcher) Run() error {
	for _, device := range deviceR.All() {
		if device.Origin {
			continue
		}

		gracePeriod := device.GracePeriodDuration()

		logger := self.logger.With(slog.Group(
			"device",
			"id", device.ID,
			"name", device.Name,
		))

		threshold := time.Now().Add(-gracePeriod)

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

		if device.State == deviceR.StateOnline && device.StateUpdatedAt.Before(time.Now().Add(-limit)) {
			if self.limitNotificationCache[device.ID] == device.StateUpdatedAt {
				continue
			}

			self.limitNotificationCache[device.ID] = device.StateUpdatedAt

			if db.B[m.DeviceLimitNotification]("device_id = ? and processed = 0", device.ID).Exists() {
				continue
			}

			db.Create(&m.DeviceLimitNotification{
				DeviceID:       device.ID,
				StateUpdatedAt: device.StateUpdatedAt,
				Limit:          device.Limit,
			})

			logger.Info("device overstayed its welcome", "limit", limit)
		}
	}

	return nil
}

func deviceOnline(device *m.Device) {
	state := deviceR.StateOnline

	db.Create(&m.DeviceStateLog{
		DeviceID:    device.ID,
		State:       state,
		GracePeriod: device.GracePeriod,
	})

	if device.Watch {
		db.Create(&m.DeviceStateNotification{
			DeviceID:    device.ID,
			State:       state,
			GracePeriod: device.GracePeriod,
		})
	}

	device.UpdateState(state)
}

func deviceOffline(device *m.Device, transitionedAt time.Time) {
	state := deviceR.StateOffline

	db.Create(&m.DeviceStateLog{
		DeviceID:    device.ID,
		State:       state,
		GracePeriod: device.GracePeriod,
		CreatedAt:   transitionedAt,
		UpdatedAt:   transitionedAt,
	})

	if device.Watch {
		db.Create(&m.DeviceStateNotification{
			DeviceID:    device.ID,
			State:       state,
			GracePeriod: device.GracePeriod,
			CreatedAt:   transitionedAt,
			UpdatedAt:   transitionedAt,
		})
	}

	device.UpdateState(state)
}
