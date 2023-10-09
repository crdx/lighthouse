package m

import (
	"time"

	"crdx.org/db"
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
	db.For[Adapter](self.ID).Update(values...)
}

// Fresh returns an Adapter with the latest values from the db.
func (self *Adapter) Fresh() *Adapter {
	i, _ := db.First[Adapter](self.ID)
	return i
}

// Delete deletes this Adapter.
func (self *Adapter) Delete() {
	db.For[Adapter](self.ID).Delete()
}

// Device returns the Device for this Adapter, and true if it has one associated.
func (self *Adapter) Device() (*Device, bool) {
	return db.First[Device](self.DeviceID)
}
