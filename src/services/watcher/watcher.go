// Package watcher watches for devices whose state has changed.
package watcher

import (
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/models/deviceM"
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
		gracePeriod := time.Duration(int64(device.GracePeriod)) * time.Second

		if device.LastSeen.Before(time.Now().Add(-gracePeriod)) {
			if device.State == deviceM.StateOnline {
				newState := deviceM.StateOffline

				db.Create(&m.DeviceStateLog{
					DeviceID: device.ID,
					State:    newState,
				})

				device.Update("state", newState)
			}
		} else {
			if device.State == deviceM.StateOffline {
				newState := deviceM.StateOnline

				db.Create(&m.DeviceStateLog{
					DeviceID: device.ID,
					State:    newState,
				})

				device.Update("state", newState)
			}
		}
	}
	return nil
}
