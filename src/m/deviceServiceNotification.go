package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceServiceNotification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID  uint `gorm:"not null;index"`
	ServiceID uint `gorm:""`
	Processed bool `gorm:"not null;default:false"`
}

func (self *DeviceServiceNotification) Update(values ...any) {
	db.For[DeviceServiceNotification](self.ID).Update(values...)
}

func (self *DeviceServiceNotification) Delete() {
	db.For[DeviceServiceNotification](self.ID).Delete()
}
