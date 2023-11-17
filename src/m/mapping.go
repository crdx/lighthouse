package m

import (
	"fmt"
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Mapping struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	MACAddress string `gorm:"size:17;not null"`
	IPAddress  string `gorm:"size:15;not null"`
	Label      string `gorm:"size:20;not null"`
}

func (self *Mapping) Update(values ...any) {
	db.For[Mapping](self.ID).Update(values...)
}

func (self *Mapping) Delete() {
	db.For[Mapping](self.ID).Delete()
}

func (self *Mapping) AuditName() string {
	return fmt.Sprintf("%s ↔ %s (%s)", self.MACAddress, self.IPAddress, self.Label)
}
