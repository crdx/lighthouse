package seeder

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceR"
)

func createDevice(id uint, name, hostname string, lastSeen time.Time) *m.Device {
	return db.Save(&m.Device{
		ID:       id,
		Name:     name,
		Hostname: hostname,
		State:    deviceR.StateOnline,
		Icon:     constants.DefaultDeviceIconClass,
		LastSeen: lastSeen,
	})
}

func createAdapter(id, deviceID uint, name, vendor, macAddress, ipAddress string, lastSeen time.Time) {
	db.Save(&m.Adapter{
		ID:         id,
		DeviceID:   deviceID,
		Name:       name,
		Vendor:     vendor,
		MACAddress: macAddress,
		IPAddress:  ipAddress,
		LastSeen:   lastSeen,
	})
}

func createDeviceStateLog(id, deviceID uint, state string, createdAt time.Time) {
	db.Save(&m.DeviceStateLog{
		ID:          id,
		DeviceID:    deviceID,
		State:       state,
		CreatedAt:   createdAt,
		GracePeriod: 5,
	})
}

func Run() error {
	t2 := time.Date(2023, time.September, 1, 12, 00, 00, 0, time.UTC)
	t1 := time.Date(2023, time.October, 1, 12, 00, 00, 0, time.UTC)
	t3 := time.Date(2023, time.November, 1, 12, 00, 00, 0, time.UTC)

	device1 := createDevice(1, "device1", "device1", t1)
	createAdapter(1, device1.ID, "adapter1", "Corp 1", "AA:AA:AA:AA:AA:AA", "127.0.0.1", t1)

	device2 := createDevice(2, "device2", "device2", t2)
	createAdapter(2, device2.ID, "adapter2", "Corp 2", "BB:BB:BB:BB:BB:BB", "127.0.0.2", t2)

	device3 := createDevice(3, "device3", "device3", t3)
	createAdapter(3, device3.ID, "adapter3", "Corp 3", "CC:CC:CC:CC:CC:CC", "127.0.0.3", t3)

	createDeviceStateLog(1, device1.ID, deviceR.StateOnline, time.Now().Add(-3*time.Minute))
	createDeviceStateLog(2, device1.ID, deviceR.StateOffline, time.Now().Add(-2*time.Minute))
	createDeviceStateLog(3, device2.ID, deviceR.StateOffline, time.Now().Add(-1*time.Minute))

	db.Save(&m.Notification{
		Subject: "a thing has happened",
		Body:    "here are more details about the thing that happened",
	})

	db.Save(&m.Setting{Name: "notification_from_header", Value: "foo <foo@example.com>"})
	db.Save(&m.Setting{Name: "notification_to_header", Value: "bar <bar@example.com>"})

	return nil
}
