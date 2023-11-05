package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceLimitNotification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID       uint      `gorm:"not null;index"`
	StateUpdatedAt time.Time `gorm:"not null"`
	Limit          uint      `gorm:"not null"`
	Processed      bool      `gorm:"not null;default:false"`
}

func (self *DeviceLimitNotification) Update(values ...any) {
	db.For[DeviceLimitNotification](self.ID).Update(values...)
}

func (self *DeviceLimitNotification) Delete() {
	db.For[DeviceLimitNotification](self.ID).Delete()
}
