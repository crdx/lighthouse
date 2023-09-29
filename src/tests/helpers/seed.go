package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Seed() {
	device := db.Save(&m.Device{MACAddress: "AA:BB:CC:DD:EE:FF", Name: "localhost"})
	db.Save(&m.DeviceMapping{DeviceID: device.ID, IPAddress: "127.0.0.1"})
}
