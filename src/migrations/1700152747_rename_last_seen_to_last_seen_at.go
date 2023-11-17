package migrations

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/util"
	"gorm.io/gorm"
)

func RenameLastSeenToLastSeenAt(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypePre,
		Run: func(db *gorm.DB) error {
			return util.Chain(
				func() error { return db.Migrator().RenameColumn(m.Device{}, "last_seen", "last_seen_at") },
				func() error { return db.Migrator().RenameColumn(m.Adapter{}, "last_seen", "last_seen_at") },
			)
		},
	}
}
