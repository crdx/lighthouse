package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Setting struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name  string `gorm:"size:100;not null"`
	Value string `gorm:"not null"`
}

func (self *Setting) Update(values ...any) {
	db.For[Setting](self.ID).Update(values...)
}

func (self *Setting) Delete() {
	db.For[Setting](self.ID).Delete()
}
