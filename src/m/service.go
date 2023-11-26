package m

import (
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/probe"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Service struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name       string    `gorm:""`
	DeviceID   uint      `gorm:"not null"`
	Port       uint      `gorm:"not null"`
	LastSeenAt time.Time `gorm:"not null"`
}

func (self *Service) Update(values ...any) {
	db.For[Service](self.ID).Update(values...)
}

func (self *Service) Delete() {
	db.For[Service](self.ID).Delete()
}

func (self *Service) Fresh() *Service {
	return lo.Must(db.First[Service](self.ID))
}

// Device returns the Device for this Service.
func (self *Service) Device() *Device {
	device, _ := db.First[Device](self.DeviceID)
	return device
}

func (self *Service) DisplayName() string {
	if self.Name == "" {
		return probe.ServiceName(self.Port)
	}
	return self.Name
}

func (self *Service) AuditName() string {
	return fmt.Sprintf("%s (ID: %d) of device %s", self.DisplayName(), self.ID, self.Device().AuditName())
}
