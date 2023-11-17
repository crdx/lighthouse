package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func MigrateToSimpleRoles(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypePost,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`UPDATE users SET role = 3 where admin = 1`).Error,
					db.Exec(`UPDATE users SET role = 1 where admin = 0`).Error,
					db.Exec(`ALTER TABLE users DROP COLUMN admin`).Error,
				)
			})
		},
	}
}
