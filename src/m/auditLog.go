package m

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/util/sqlutil"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type AuditLog struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID    sqlutil.NullUint `gorm:""`
	DeviceID  sqlutil.NullUint `gorm:""`
	IPAddress string           `gorm:"size:15;not null"`
	Message   string           `gorm:"not null"`
}

func (self *AuditLog) Update(values ...any) {
	db.For[AuditLog](self.ID).Update(values...)
}

func (self *AuditLog) Delete() {
	db.For[AuditLog](self.ID).Delete()
}
