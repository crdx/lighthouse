package conf

import (
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
}
