package m

import (
	"database/sql"
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Scan struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CompletedAt sql.NullTime `gorm:""`
}

func (self *Scan) Update(values ...any) {
	db.For[Scan](self.ID).Update(values...)
}

func (self *Scan) Delete() {
	db.For[Scan](self.ID).Delete()
}
