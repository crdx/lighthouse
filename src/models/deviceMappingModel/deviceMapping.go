package deviceMappingModel

import (
	"time"

	"crdx.org/db"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type DeviceMapping struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:""`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeviceID  uint           `gorm:"not null"`
	IPAddress string         `gorm:"size:15;not null"`
	LastSeen  time.Time      `gorm:"not null"`
}

func Upsert(deviceID uint, ipAddress string) (DeviceMapping, bool) {
	// This method has the potential to be called very often, so let's not hang onto a model object
	// for longer than necessary. Immediately create the record if it doesn't exist, and then just
	// run the query that updates the specific fields.
	deviceMapping, found := db.FirstOrCreate(DeviceMapping{
		DeviceID:  deviceID,
		IPAddress: ipAddress,
	})

	columns := db.Map{
		"last_seen": time.Now(),
	}

	q := DeviceMapping{ID: deviceMapping.ID}

	db.B(q).Update(columns)
	deviceMapping, _ = db.B(q).First()

	return deviceMapping, found
}
