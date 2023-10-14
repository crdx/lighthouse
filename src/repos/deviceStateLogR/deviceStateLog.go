package deviceStateLogR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func All() []*m.DeviceStateLog {
	return db.B[m.DeviceStateLog]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func LatestForDevice(deviceID uint) []*m.DeviceStateLog {
	return db.B(m.DeviceStateLog{DeviceID: deviceID}).Limit(5).Order("created_at DESC").Find()
}
