package m

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/pkg/duration"
	"crdx.org/lighthouse/pkg/util"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Device struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Origin              bool         `gorm:"not null"`
	Name                string       `gorm:"size:100;not null"`
	Hostname            string       `gorm:"size:100;not null"`
	HostnameAnnouncedAt sql.NullTime `gorm:""`
	State               string       `gorm:"size:20;not null"`
	StateUpdatedAt      time.Time    `gorm:"not null"`
	Icon                string       `gorm:"size:100;not null"`
	Notes               string       `gorm:"size:5000,not null"`
	LastSeenAt          time.Time    `gorm:"not null"`
	Watch               bool         `gorm:"not null"`
	Ping                bool         `gorm:"not null"`
	Limit               string       `gorm:"not null"`
	GracePeriod         string       `gorm:"not null"`
}

func (self *Device) Update(values ...any) {
	db.For[Device](self.ID).Update(values...)
}

func (self *Device) Delete() {
	for _, adapter := range self.Adapters() {
		adapter.Delete()
	}

	db.B[DeviceStateLog]("device_id = ?", self.ID).Delete()
	db.B[DeviceStateNotification]("device_id = ?", self.ID).Delete()
	db.B[DeviceDiscoveryNotification]("device_id = ?", self.ID).Delete()

	db.For[Device](self.ID).Delete()
}

func (self *Device) Fresh() *Device {
	return lo.Must(db.First[Device](self.ID))
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

	if self.Name == "" || db.B[Device]("name = ?", self.Name).Count() > 1 {
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
	return db.B[Adapter]("device_id = ?", self.ID).Order("last_seen_at DESC").Find()
}

// Services returns all Services found for this Device.
func (self *Device) Services() []*Service {
	return db.B[Service]("device_id = ?", self.ID).Order("port ASC").Find()
}

func (self *Device) UpdateState(state string) {
	self.Update(
		"state", state,
		"state_updated_at", time.Now(),
	)
}

func (self *Device) IconClass() string {
	return util.IconToClass(self.Icon)
}

func (self *Device) LimitDuration() time.Duration {
	return lo.Must(duration.Parse(self.Limit))
}

func (self *Device) GracePeriodDuration() time.Duration {
	return lo.Must(duration.Parse(self.GracePeriod))
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
