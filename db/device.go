package db

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/duration"
	"crdx.org/lighthouse/pkg/util"
	"crdx.org/lighthouse/pkg/util/timeutil"
)

func (self *Device) CascadeDelete() {
	for _, adapter := range self.Adapters() {
		adapter.Delete()
	}

	if v, found := FindDeviceStateLogByDeviceID(self.ID); found {
		v.Delete()
	}
	if v, found := FindDeviceStateNotificationByDeviceID(self.ID); found {
		v.Delete()
	}
	if v, found := FindDeviceDiscoveryNotificationByDeviceID(self.ID); found {
		v.Delete()
	}

	self.Delete()
}

func (self *Device) DisplayName() string {
	if self.Name == "" {
		return "Unknown"
	} else {
		return self.Name
	}
}

// Identifier returns the name of this device and if it's not uniquely named then appends the
// hostname (if it exists) or the ID.
func (self *Device) Identifier() string {
	id := self.DisplayName()

	if self.Name == "" || len(FindDevicesByName(self.Name)) > 1 {
		if self.Hostname != "" {
			id += fmt.Sprintf(" (%s)", self.Hostname)
		} else {
			id += fmt.Sprintf(" (%d)", self.ID)
		}
	}

	return id
}

func (self *Device) Details() string {
	var s strings.Builder

	discovered := timeutil.ToLocal(self.CreatedAt).Format(constants.TimeFormatReadablePrecise)

	// Ensure HTML minification can't strip these spaces.
	indent := "\u00A0\u00A0\u00A0\u00A0"

	s.WriteString(fmt.Sprintf("%s:\n", self.DisplayName()))
	s.WriteString(fmt.Sprintf(indent+"Discovered: %s\n", discovered))
	if self.Hostname != "" {
		s.WriteString(fmt.Sprintf(indent+"Hostname: %s\n", self.Hostname))
	}

	for _, adapter := range self.Adapters() {
		s.WriteString(fmt.Sprintf(indent+"MAC Address: %s\n", adapter.MACAddress))
		s.WriteString(fmt.Sprintf(indent+"IP Address: %s\n", adapter.IPAddress))
		if adapter.Vendor != "" {
			s.WriteString(fmt.Sprintf(indent+"Vendor: %s\n", adapter.Vendor))
		}
	}

	return strings.TrimSpace(s.String())
}

// Adapters returns all Adapters attached to this Device.
func (self *Device) Adapters() []*Adapter {
	return FindAdaptersForDevice(self.ID)
}

// Services returns all Services found for this Device.
func (self *Device) Services() []*Service {
	return FindServicesForDevice(self.ID)
}

func (self *Device) IconClass() string {
	return util.IconToClass(self.Icon)
}

func (self *Device) LimitDuration() time.Duration {
	return duration.MustParse(self.Limit)
}

func (self *Device) GracePeriodDuration() time.Duration {
	return duration.MustParse(self.GracePeriod)
}

func (self *Device) AuditName() string {
	return fmt.Sprintf("%s (ID: %d)", self.DisplayName(), self.ID)
}

func (self *Device) LogAttr() slog.Attr {
	return slog.Group(
		"device",
		"id", self.ID,
		"name", self.DisplayName(),
		"hostname", self.Hostname,
	)
}
