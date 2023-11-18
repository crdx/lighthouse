package m

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Notification struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Subject string `gorm:"size:200;not null"`
	Body    string `gorm:"not null"`
}

func (self *Notification) Update(values ...any) {
	db.For[Notification](self.ID).Update(values...)
}

func (self *Notification) Delete() {
	db.For[Notification](self.ID).Delete()
}
