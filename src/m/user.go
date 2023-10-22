package m

import (
	"database/sql"
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type User struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username     string       `gorm:"size:20;not null"`
	PasswordHash string       `gorm:"not null"`
	Admin        bool         `gorm:"not null"`
	LastLogin    sql.NullTime `gorm:""`
}

func (self *User) Update(values ...any) {
	db.For[User](self.ID).Update(values...)
}

func (self *User) Delete() {
	db.For[User](self.ID).Delete()
}
