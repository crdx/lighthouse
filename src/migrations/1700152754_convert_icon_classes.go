package migrations

import (
	"regexp"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"gorm.io/gorm"
)

func ConvertIconClasses(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypePost,
		Run: func(db *gorm.DB) error {
			var devices []m.Device
			db.Find(&devices)
			for _, device := range devices {
				if device.Icon != "" {
					parts := regexp.MustCompile(`fa-(\S+) fa-(\S+)`).FindStringSubmatch(device.Icon) //nolint
					device.Icon = parts[1] + ":" + parts[2]
					db.Save(&device)
				}
			}
			return nil
		},
	}
}
