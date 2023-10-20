package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type VendorLookup struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AdapterID uint           `gorm:"not null"`
	Processed bool           `gorm:"not null;index"`
	Succeeded bool           `gorm:"not null"`
}

func (self *VendorLookup) Update(values ...any) {
	db.For[VendorLookup](self.ID).Update(values...)
}

func (self *VendorLookup) Delete() {
	db.For[VendorLookup](self.ID).Delete()
}
