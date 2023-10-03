package deviceStateLogM

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
	DeviceID  uint           `gorm:"not null"`
	State     string         `gorm:"size:15;not null"`
}

func (self *DeviceStateLog) Update(values ...any) {
	For(self.ID).Update(values...)
}

func For(id uint) *db.Builder[DeviceStateLog] {
	return db.B(DeviceStateLog{ID: id})
}
