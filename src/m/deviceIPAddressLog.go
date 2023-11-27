package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceIPAddressLog struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DeviceID  uint   `gorm:"not null;index"`
	IPAddress string `gorm:"not null"`
}

func (self *DeviceIPAddressLog) Update(values ...any) {
	db.For[DeviceIPAddressLog](self.ID).Update(values...)
}

func (self *DeviceIPAddressLog) Delete() {
	db.For[DeviceIPAddressLog](self.ID).Delete()
}
