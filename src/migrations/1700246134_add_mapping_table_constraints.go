package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddMappingTableConstraints(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`ALTER TABLE mappings ADD CONSTRAINT chk_mappings__mac_address_valid CHECK (mac_address REGEXP '^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$')`).Error,
					db.Exec(`ALTER TABLE mappings ADD CONSTRAINT chk_mappings__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,
				)
			})
		},
	}
}
