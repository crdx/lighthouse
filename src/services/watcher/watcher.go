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
	log                    *slog.Logger
	limitNotificationCache map[uint]time.Time
}

func New() *Watcher {
	return &Watcher{}
}

func (self *Watcher) Init(args *services.Args) error {
	self.log = args.Logger
	self.limitNotificationCache = map[uint]time.Time{}

	return nil
}

func (self *Watcher) Run() error {
	for _, device := range deviceR.All() {
		gracePeriod := time.Duration(int64(device.GracePeriod)) * time.Minute

		log := self.log.With(slog.Group(
			"device",
			"id", device.ID,
			"name", device.Name,
		))

		if device.Origin {
			continue
		}

		if device.LastSeenAt.Before(time.Now().Add(-gracePeriod)) {
			if device.State == deviceR.StateOnline {
				deviceOffline(device)
				log.Info("device is offline")
			}
		} else {
			if device.State == deviceR.StateOffline {
				deviceOnline(device)
				log.Info("device is online")
			}
		}

		if device.Limit == 0 {
			continue
		}

		limit := time.Duration(int64(device.Limit)) * time.Minute

		if device.State == deviceR.StateOnline && device.StateUpdatedAt.Before(time.Now().Add(-limit)) {
			if self.limitNotificationCache[device.ID] == device.StateUpdatedAt {
				continue
			}

			if db.B[m.DeviceLimitNotification]("device_id = ? and processed = 0", device.ID).Exists() {
				continue
			}

			self.limitNotificationCache[device.ID] = device.StateUpdatedAt

			db.Create(&m.DeviceLimitNotification{
				DeviceID:       device.ID,
				StateUpdatedAt: device.StateUpdatedAt,
				Limit:          device.Limit,
			})

			log.Info("device overstayed its welcome", "limit", limit)
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

func deviceOffline(device *m.Device) {
	state := deviceR.StateOffline

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
