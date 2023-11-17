package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddForeignKeyConstraints(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return errors.Join(
				db.Exec(`ALTER TABLE audit_logs CHANGE COLUMN user_id user_id BIGINT UNSIGNED NULL`).Error, //nolint
				db.Exec(`UPDATE audit_logs SET user_id = NULL where user_id = 0`).Error,
				db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT fk_audit_logs__user FOREIGN KEY (user_id) REFERENCES users (id)`).Error,

				db.Exec(`ALTER TABLE adapters ADD CONSTRAINT fk_adapters__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				db.Exec(`ALTER TABLE device_discovery_notifications ADD CONSTRAINT fk_device_discovery_notifications__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				db.Exec(`ALTER TABLE device_limit_notifications ADD CONSTRAINT fk_device_limit_notifications__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				db.Exec(`ALTER TABLE device_state_logs ADD CONSTRAINT fk_device_state_logs__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				db.Exec(`ALTER TABLE device_state_notifications ADD CONSTRAINT fk_device_state_notifications__device FOREIGN KEY (device_id) REFERENCES devices (id)`).Error,
				db.Exec(`ALTER TABLE vendor_lookups ADD CONSTRAINT fk_vendor_lookups__adapter FOREIGN KEY (adapter_id) REFERENCES adapters (id)`).Error,
			)
		},
	}
}
