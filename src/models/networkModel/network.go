package networkModel

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Network struct {
	ID                uint           `gorm:"primarykey"`
	CreatedAt         time.Time      `gorm:""`
	UpdatedAt         time.Time      `gorm:""`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	IPRange           string         `gorm:"size:18;not null"`
	Name              string         `gorm:"size:255;not null"`
	AlertOnNewDevices bool           `gorm:"default:false"`
}

func Upsert(ipRange string) (Network, bool) {
	network, found := db.FirstOrInit(Network{IPRange: ipRange})

	if !found {
		network.Name = ipRange
		network.AlertOnNewDevices = true

		db.Save(&network)
	}

	return network, found
}
