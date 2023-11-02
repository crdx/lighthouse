package conf

import (
	"crdx.org/db"
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
}
