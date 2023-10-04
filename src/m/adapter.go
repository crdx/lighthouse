package m

import (
	"time"

	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Adapter struct {
	ID         uint           `gorm:"primarykey"`
	CreatedAt  time.Time      `gorm:""`
	UpdatedAt  time.Time      `gorm:""`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	DeviceID   uint           `gorm:"not null"`
	Name       string         `gorm:"size:255;not null"`
	MACAddress string         `gorm:"size:17;not null"`
	Vendor     string         `gorm:"size:255;not null"`
	IPAddress  string         `gorm:"size:15;not null"`
	LastSeen   time.Time      `gorm:"not null"`
}

func (self *Adapter) Update(values ...any) {
	ForAdapter(self.ID).Update(values...)
}

// Fresh returns an Adapter with the latest values from the db.
func (self *Adapter) Fresh() *Adapter {
	i, _ := ForAdapter(self.ID).First()
	return i
}

// Delete deletes this Adapter.
func (self *Adapter) Delete() {
	ForAdapter(self.ID).Delete()
}

// Device returns the Device for this Adapter, and true if it has one associated.
func (self *Adapter) Device() (*Device, bool) {
	return ForDevice(self.DeviceID).First()
}
