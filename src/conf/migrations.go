package conf

import (
	"errors"
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
	{
		ID:   "ConvertDurations",
		Type: db.MigrationTypePost,
		Run: func(db *gorm.DB) error {
			return errors.Join(
				db.Exec(`UPDATE devices SET grace_period = CONCAT(grace_period, ' mins')`).Error,
				db.Exec(`UPDATE device_state_logs SET grace_period = CONCAT(grace_period, ' mins')`).Error,
				db.Exec(`UPDATE device_state_notifications SET grace_period = CONCAT(grace_period, ' mins')`).Error,
				db.Exec(`UPDATE settings SET value = CONCAT(value, ' min') where name = 'scan_interval'`).Error,
			)
		},
	},
}
