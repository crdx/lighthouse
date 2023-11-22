package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddAuditLogDeviceIdConstraint(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT fk_audit_logs__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				)
			})
		},
	}
}
