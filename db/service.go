package db

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/pkg/probe"
)

func (self *Service) DisplayName() string {
	if self.Name == "" {
		return probe.ServiceName(self.Port)
	}
	return self.Name
}

func (self *Service) Device() *Device {
	device, _ := FindDevice(self.DeviceID)
	return device
}

func (self *Service) AuditName() string {
	return fmt.Sprintf("%s (ID: %d) of device %s", self.DisplayName(), self.ID, self.Device().AuditName())
}

func (self *Service) Details() string {
	var s strings.Builder

	name := self.DisplayName()
	if name == "" {
		name = "unknown"
	}

	s.WriteString(fmt.Sprintf("port %d (%s) is open on %s", self.Port, name, self.Device().DisplayName()))
	return strings.TrimSpace(s.String())
}
