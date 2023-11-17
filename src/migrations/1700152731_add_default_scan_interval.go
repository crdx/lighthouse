package migrations

import (
	"crdx.org/db"
	"gorm.io/gorm"
)

func AddDefaultScanInterval(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypePre,
		Run: func(db *gorm.DB) error {
			return db.Exec(`INSERT INTO settings (name, value) VALUES ('scan_interval', '1')`).Error
		},
	}
}
