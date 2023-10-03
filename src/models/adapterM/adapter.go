package adapterM

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/util"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Adapter struct {
	ID         uint           `gorm:"primarykey"`
	CreatedAt  time.Time      `gorm:""`
	UpdatedAt  time.Time      `gorm:""`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	DeviceID   uint           `gorm:"not null"`
	Name       string         `gorm:"size:255;not null"`
	MACAddress string         `gorm:"size:17;not null"`
	Vendor     string         `gorm:"size:255;not null"`
	IPAddress  string         `gorm:"size:15;not null"`
	LastSeen   time.Time      `gorm:"not null"`
}

func (self *Adapter) Update(values ...any) {
	For(self.ID).Update(values...)
}

func (self *Adapter) Fresh() *Adapter {
	i, _ := For(self.ID).First()
	return i
}

func (self *Adapter) Delete() {
	For(self.ID).Delete()
}

func For(id uint) *db.Builder[Adapter] {
	return db.B(Adapter{ID: id})
}

func Upsert(macAddress string, ipAddress string) (*Adapter, bool) {
	adapter, adapterFound := db.FirstOrCreate(Adapter{
		MACAddress: macAddress,
	})

	columns := db.Map{
		"last_seen":  time.Now(),
		"ip_address": ipAddress,
	}

	if adapter.Vendor == "" {
		if vendor, vendorFound := util.GetVendor(adapter.MACAddress); vendorFound {
			columns["vendor"] = vendor

			if !adapterFound {
				columns["name"] = vendor
			}
		}
	}

	adapter.Update(columns)

	return adapter.Fresh(), adapterFound
}
