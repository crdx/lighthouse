package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Seed() {
	device := db.Save(&m.Device{Name: "localhost"})
	db.Save(&m.Adapter{DeviceID: device.ID, MACAddress: "AA:BB:CC:DD:EE:FF", IPAddress: "127.0.0.1"})
}
