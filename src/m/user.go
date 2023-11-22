package m

import (
	"database/sql"
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"github.com/samber/lo"
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
	Role         uint         `gorm:"not null"`
	LastLoginAt  sql.NullTime `gorm:""`
	LastVisitAt  sql.NullTime `gorm:""`
}

func (self *User) Update(values ...any) {
	db.For[User](self.ID).Update(values...)
}

func (self *User) Delete() {
	db.For[User](self.ID).Delete()
}

func (self *User) Fresh() *User {
	return lo.Must(db.First[User](self.ID))
}

func (self *User) AuditName() string {
	return fmt.Sprintf("%s (ID: %d)", self.Username, self.ID)
}

func (self *User) IsAdmin() bool  { return self.Role >= constants.RoleAdmin }
func (self *User) IsEditor() bool { return self.Role >= constants.RoleEditor }
func (self *User) IsViewer() bool { return self.Role >= constants.RoleViewer }
