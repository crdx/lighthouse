// Package watcher watches for devices whose state has changed.
package watcher

import (
	"time"

	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceR"
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
	for _, device := range deviceR.All() {
		gracePeriod := time.Duration(int64(device.GracePeriod)) * time.Minute

		if device.LastSeen.Before(time.Now().Add(-gracePeriod)) {
			if device.State == deviceR.StateOnline {
				newState := deviceR.StateOffline

				db.Create(&m.DeviceStateLog{
					DeviceID: device.ID,
					State:    newState,
				})

				device.Update("state", newState)
			}
		} else {
			if device.State == deviceR.StateOffline {
				newState := deviceR.StateOnline

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
