package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type ScanResult struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ScanID   uint `gorm:"not null"`
	DeviceID uint `gorm:"not null"`
	Port     uint `gorm:"not null"`
}

func (self *ScanResult) Update(values ...any) {
	db.For[ScanResult](self.ID).Update(values...)
}

func (self *ScanResult) Delete() {
	db.For[ScanResult](self.ID).Delete()
}
