// Package watcher watches for devices whose state has changed.
package watcher

import (
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
	"crdx.org/lighthouse/services"
)

type Watcher struct {
	log *slog.Logger
}

func New() *Watcher {
	return &Watcher{}
}

func (self *Watcher) Init(args *services.Args) error {
	self.log = args.Logger
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
