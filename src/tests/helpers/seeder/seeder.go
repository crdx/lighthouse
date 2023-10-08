package seeder

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceR"
)

func createNetwork(id uint, name, ipRange string) *m.Network {
	return db.Save(&m.Network{
		ID:      id,
		Name:    name,
		IPRange: ipRange,
	})
}

func createDevice(id, networkID uint, name, hostname string, lastSeen time.Time) *m.Device {
	return db.Save(&m.Device{
		ID:        id,
		NetworkID: networkID,
		Name:      name,
		Hostname:  hostname,
		State:     deviceR.StateOnline,
		Icon:      constants.DefaultDeviceIconClass,
		LastSeen:  lastSeen,
	})
}

func createAdapter(id, deviceID uint, name, vendor, macAddress, ipAddress string, lastSeen time.Time) *m.Adapter {
	return db.Save(&m.Adapter{
		ID:         id,
		DeviceID:   deviceID,
		Name:       name,
		Vendor:     vendor,
		MACAddress: macAddress,
		IPAddress:  ipAddress,
		LastSeen:   lastSeen,
	})
}

func Run() {
	t1 := time.Date(2023, time.October, 1, 12, 00, 00, 0, time.UTC)
	t2 := time.Date(2023, time.September, 1, 12, 00, 00, 0, time.UTC)

	network1 := createNetwork(1, "network1", "127.0.0.1/24")

	device1 := createDevice(1, network1.ID, "device1", "device1", t1)
	createAdapter(1, device1.ID, "adapter1", "Corp 1", "AA:AA:AA:AA:AA:AA", "127.0.0.1", t1)

	device2 := createDevice(2, network1.ID, "device2", "device2", t2)
	createAdapter(2, device2.ID, "adapter2", "Corp 2", "BB:BB:BB:BB:BB:BB", "127.0.0.2", t2)
}