package helpers

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceR"
)

func Seed() {
	device := db.Save(&m.Device{
		ID:       1,
		Name:     "localhost",
		Hostname: "localhost",
		State:    deviceR.StateOnline,
		Icon:     constants.DefaultDeviceIconClass,
		LastSeen: time.Date(2000, time.January, 1, 12, 00, 00, 0, time.UTC),
	})

	db.Save(&m.Adapter{
		ID:         1,
		DeviceID:   device.ID,
		Name:       "adapter1",
		Vendor:     "Computer Corporation",
		MACAddress: "AA:BB:CC:DD:EE:FF",
		IPAddress:  "127.0.0.1",
		LastSeen:   time.Date(2000, time.January, 1, 12, 00, 00, 0, time.UTC),
	})
}
