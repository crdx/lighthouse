package db

import (
	"fmt"
	"log/slog"
)

// Device returns the Device for this Adapter.
func (self *Adapter) Device() *Device {
	device, _ := FindDevice(self.DeviceID)
	return device
}

func (self *Adapter) AuditName() string {
	return fmt.Sprintf("%s (ID: %d) of device %s", self.Name, self.ID, self.Device().AuditName())
}

func (self *Adapter) LogAttr() slog.Attr {
	return slog.Group(
		"adapter",
		"id", self.ID,
		"mac", self.MACAddress,
		"ip", self.IPAddress,
		"vendor", self.Vendor,
	)
}
