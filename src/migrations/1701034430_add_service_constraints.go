package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddServiceConstraints(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`ALTER TABLE services ADD CONSTRAINT fk_services__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
					db.Exec(`ALTER TABLE device_service_notifications ADD CONSTRAINT fk_device_service_notifications__service FOREIGN KEY (service_id) REFERENCES services (id)`).Error,
				)
			})
		},
	}
}
