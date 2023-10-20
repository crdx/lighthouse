package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceStateNotification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID    uint   `gorm:"not null;index"`
	State       string `gorm:"size:32;not null"`
	GracePeriod uint   `gorm:"not null"`
	Processed   bool   `gorm:"not null;default:false"`
}

func (self *DeviceStateNotification) Update(values ...any) {
	db.For[DeviceStateNotification](self.ID).Update(values...)
}

func (self *DeviceStateNotification) Delete() {
	db.For[DeviceStateNotification](self.ID).Delete()
}
