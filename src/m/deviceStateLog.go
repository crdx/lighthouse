package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceStateLog struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID    uint   `gorm:"not null"`
	State       string `gorm:"size:20;not null"`
	GracePeriod string `gorm:"not null"`
}

func (self *DeviceStateLog) Update(values ...any) {
	db.For[DeviceStateLog](self.ID).Update(values...)
}

func (self *DeviceStateLog) Delete() {
	db.For[DeviceStateLog](self.ID).Delete()
}
