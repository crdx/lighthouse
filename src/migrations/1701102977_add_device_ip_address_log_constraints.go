package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddDeviceIpAddressLogConstraints(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`ALTER TABLE device_ip_address_logs ADD CONSTRAINT fk_device_ip_address_logs__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
					db.Exec(`ALTER TABLE device_ip_address_logs ADD CONSTRAINT chk_device_ip_address_logs__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,
				)
			})
		},
	}
}
