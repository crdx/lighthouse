// Package watcher watches for devices whose state has changed.
package watcher

import (
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/models/deviceModel"
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

func (*Watcher) Run() error {
	for _, device := range db.B[m.Device]().Find() {
		q := m.Device{ID: device.ID}

		gracePeriod := time.Duration(int64(device.GracePeriod)) * time.Second

		if device.LastSeen.Before(time.Now().Add(-gracePeriod)) {
			if device.State == deviceModel.StateOnline {
				newState := deviceModel.StateOffline

				db.Create(&m.DeviceStateLog{
					DeviceID: device.ID,
					State:    newState,
				})

				db.B(q).Update("state", newState)
			}
		} else {
			if device.State == deviceModel.StateOffline {
				newState := deviceModel.StateOnline

				db.Create(&m.DeviceStateLog{
					DeviceID: device.ID,
					State:    newState,
				})

				db.B(q).Update("state", newState)
			}
		}
	}
	return nil
}
