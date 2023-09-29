package notificationModel

import (
	"time"

	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Notification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Message   string         `gorm:"size:255;not null"`
	Processed bool           `gorm:"default:false;not null"`
}
