package conf

import (
	"regexp"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/util"
	"gorm.io/gorm"
)

var migrations = []*db.Migration{
	{
		ID:   "AddDefaultScanInterval",
		Type: db.MigrationTypePre,
		Run: func(db *gorm.DB) error {
			return db.Exec(`INSERT INTO settings (name, value) VALUES ('scan_interval', '1')`).Error
		},
	},
	{
		ID:   "RenameLastSeenToLastSeenAt",
		Type: db.MigrationTypePre,
		Run: func(db *gorm.DB) error {
			return util.Chain(
				func() error { return db.Migrator().RenameColumn(m.Device{}, "last_seen", "last_seen_at") },
				func() error { return db.Migrator().RenameColumn(m.Adapter{}, "last_seen", "last_seen_at") },
			)
		},
	},
	{
		ID:   "ConvertIconClasses",
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
	},
}
