package m

import (
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
	ForDevice(self.ID).Update(values...)
}

// Fresh returns a Device with the latest values from the db.
func (self *Device) Fresh() *Device {
	i, _ := ForDevice(self.ID).First()
	return i
}

// Delete deletes this Device and all of its attached Adapters.
func (self *Device) Delete() {
	for _, adapter := range self.Adapters() {
		adapter.Delete()
	}

	ForDevice(self.ID).Delete()
}

// Adapters returns all Adapters attached to this Device.
func (self *Device) Adapters() []*Adapter {
	return db.B(Adapter{DeviceID: self.ID}).Find()
}
