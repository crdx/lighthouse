package m

import (
	"fmt"
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Device struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `gorm:""`
	UpdatedAt   time.Time      `gorm:""`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	NetworkID   uint           `gorm:"not null"`
	Name        string         `gorm:"size:255;not null"`
	Hostname    string         `gorm:"size:255;not null"`
	State       string         `gorm:"size:32;not null"`
	Icon        string         `gorm:"size:255;not null"`
	Notes       string         `gorm:"not null"`
	LastSeen    time.Time      `gorm:""`
	Watch       bool           `gorm:"not null;default:false"`
	GracePeriod uint           `gorm:"not null;default:5"`
}

func (self *Device) Update(values ...any) {
	db.For[Device](self.ID).Update(values...)
}

func (self *Device) Delete() {
	for _, adapter := range self.Adapters() {
		adapter.Delete()
	}

	db.B[DeviceStateLog]().Where("device_id = ?", self.ID).Delete()
	db.B[DeviceStateNotification]().Where("device_id = ?", self.ID).Delete()

	db.For[Device](self.ID).Delete()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func (self *Device) DisplayName() string {
	var s string

	if self.Name == "" {
		s = "Untitled"
	} else {
		s = self.Name
	}

	// Also show the hostname (if set, otherwise ID) if this device's name is not unique across the
	// rest of the devices.
	if db.B(Device{Name: self.Name}).Count() > 1 {
		if self.Hostname != "" {
			s += fmt.Sprintf(" (%s)", self.Hostname)
		} else {
			s += fmt.Sprintf(" (%d)", self.ID)
		}
	}

	return s
}

// Adapters returns all Adapters attached to this Device.
func (self *Device) Adapters() []*Adapter {
	return db.B(Adapter{DeviceID: self.ID}).Find()
}
