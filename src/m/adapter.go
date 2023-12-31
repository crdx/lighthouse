package m

import (
	"fmt"
	"log/slog"
	"time"

	"crdx.org/db"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Adapter struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID   uint      `gorm:"not null"`
	Name       string    `gorm:"size:100;not null"`
	MACAddress string    `gorm:"size:17;not null;index"`
	Vendor     string    `gorm:"size:100;not null"`
	IPAddress  string    `gorm:"size:15;not null"`
	LastSeenAt time.Time `gorm:"not null"`
}

func (self *Adapter) Update(values ...any) {
	db.For[Adapter](self.ID).Update(values...)
}

func (self *Adapter) Delete() {
	db.B[VendorLookup]("adapter_id = ?", self.ID).Delete()

	db.For[Adapter](self.ID).Delete()
}

func (self *Adapter) Fresh() *Adapter {
	return lo.Must(db.First[Adapter](self.ID))
}

// Device returns the Device for this Adapter.
func (self *Adapter) Device() *Device {
	device, _ := db.First[Device](self.DeviceID)
	return device
}

// IsOnline returns true if this adapter was last seen within its grace period.
func (self *Adapter) IsOnline() bool {
	return self.LastSeenAt.After(time.Now().Add(-self.Device().GracePeriodDuration()))
}

// IsNotResponding returns true if this adapter was last seen within half of the grace period. This
// indicates a device that may be about to go offline.
func (self *Adapter) IsNotResponding() bool {
	return self.LastSeenAt.Before(time.Now().Add(-self.Device().GracePeriodDuration() / 2))
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
