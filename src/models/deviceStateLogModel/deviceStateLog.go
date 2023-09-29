package deviceStateLogModel

import (
	"time"

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
