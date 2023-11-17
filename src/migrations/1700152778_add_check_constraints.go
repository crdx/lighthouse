package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func AddCheckConstraints(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return errors.Join(
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT chk_devices__icon_valid CHECK (icon REGEXP '^.+:.+$' OR icon = '')`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT chk_devices__watch_valid CHECK (watch IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT chk_devices__origin_valid CHECK (origin IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT chk_devices__grace_period_valid CHECK (grace_period != '')`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT chk_devices__state_valid CHECK (state IN ('online', 'offline'))`).Error,

				db.Exec(`ALTER TABLE adapters ADD CONSTRAINT chk_adapters__mac_address_valid CHECK (mac_address REGEXP '^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$')`).Error,
				db.Exec(`ALTER TABLE adapters ADD CONSTRAINT chk_adapters__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,

				db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT chk_audit_logs__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,
				db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT chk_audit_logs__message_valid CHECK (message != '')`).Error,

				db.Exec(`ALTER TABLE device_state_logs ADD CONSTRAINT chk_device_state_logs__state_valid CHECK (state IN ('online', 'offline'))`).Error,
				db.Exec(`ALTER TABLE device_state_notifications ADD CONSTRAINT chk_device_state_notifications__state_valid CHECK (state IN ('online', 'offline'))`).Error,
				db.Exec(`ALTER TABLE device_state_notifications ADD CONSTRAINT chk_device_state_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE device_limit_notifications ADD CONSTRAINT chk_device_limit_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE device_discovery_notifications ADD CONSTRAINT chk_device_discovery_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,

				db.Exec(`ALTER TABLE vendor_lookups ADD CONSTRAINT chk_vendor_lookups__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE vendor_lookups ADD CONSTRAINT chk_vendor_lookups__succeeded_valid CHECK (succeeded IN (0, 1))`).Error,

				db.Exec(`ALTER TABLE users ADD CONSTRAINT chk_users__role_valid CHECK (role IN (1, 2, 3))`).Error,
				db.Exec(`ALTER TABLE users ADD CONSTRAINT chk_users__last_login_at_valid CHECK (last_login_at >= created_at)`).Error,
			)
		},
	}
}
