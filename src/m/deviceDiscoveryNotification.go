package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceDiscoveryNotification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeviceID  uint           `gorm:"not null;index"`
	Processed bool           `gorm:"not null;default:false"`
}

func (self *DeviceDiscoveryNotification) Update(values ...any) {
	db.For[DeviceDiscoveryNotification](self.ID).Update(values...)
}

func (self *DeviceDiscoveryNotification) Delete() {
	db.For[DeviceDiscoveryNotification](self.ID).Delete()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————
