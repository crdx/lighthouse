package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func RenameSettings(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypePre,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`UPDATE settings SET name = 'device_scan_interval' WHERE name = 'scan_interval'`).Error,
					db.Exec(`UPDATE settings SET name = 'enable_device_scan' WHERE name = 'passive'`).Error,
					db.Exec(`UPDATE settings SET name = 'notify_on_new_device' WHERE name = 'watch'`).Error,
				)
			})
		},
	}
}
