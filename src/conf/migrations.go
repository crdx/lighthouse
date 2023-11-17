package conf

import (
	"errors"
	"regexp"

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
	{
		ID:   "ConvertIconClasses",
		Type: db.MigrationTypePost,
		Run: func(db *gorm.DB) error {
			var devices []m.Device
			db.Find(&devices)
			for _, device := range devices {
				if device.Icon != "" {
					parts := regexp.MustCompile(`fa-(\S+) fa-(\S+)`).FindStringSubmatch(device.Icon) //nolint
					device.Icon = parts[1] + ":" + parts[2]
					db.Save(&device)
				}
			}
			return nil
		},
	},
	{
		ID:   "ConvertDurations",
		Type: db.MigrationTypePost,
		Run: func(db *gorm.DB) error {
			return db.Transaction(func(db *gorm.DB) error {
				return errors.Join(
					db.Exec(`UPDATE devices SET grace_period = CONCAT(grace_period, ' mins')`).Error,
					db.Exec(`UPDATE device_state_logs SET grace_period = CONCAT(grace_period, ' mins')`).Error,
					db.Exec(`UPDATE device_state_notifications SET grace_period = CONCAT(grace_period, ' mins')`).Error,
					db.Exec(`UPDATE settings SET value = CONCAT(value, ' min') where name = 'scan_interval'`).Error,
				)
			})
		},
	},
	{
		ID:   "MigrateToSimpleRoles",
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
	},
	{
		ID:   "AddForeignKeyConstraints",
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return errors.Join(
				db.Exec(`ALTER TABLE audit_logs CHANGE COLUMN user_id user_id BIGINT UNSIGNED NULL`).Error,
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
	},
	{
		ID:   "AddCheckConstraints",
		Type: db.MigrationTypeSchema,
		Run: func(db *gorm.DB) error {
			return errors.Join(
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT ck_devices__icon_valid CHECK (icon REGEXP '^.+:.+$' OR icon = '')`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT ck_devices__watch_valid CHECK (watch IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT ck_devices__origin_valid CHECK (origin IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT ck_devices__grace_period_valid CHECK (grace_period != '')`).Error,
				db.Exec(`ALTER TABLE devices ADD CONSTRAINT ck_devices__state_valid CHECK (state IN ('online', 'offline'))`).Error,

				db.Exec(`ALTER TABLE adapters ADD CONSTRAINT ck_adapters__mac_address_valid CHECK (mac_address REGEXP '^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$')`).Error,
				db.Exec(`ALTER TABLE adapters ADD CONSTRAINT ck_adapters__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,

				db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT ck_audit_logs__ip_address_valid CHECK (ip_address REGEXP '^[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}$')`).Error,
				db.Exec(`ALTER TABLE audit_logs ADD CONSTRAINT ck_audit_logs__message_valid CHECK (message != '')`).Error,

				db.Exec(`ALTER TABLE device_state_logs ADD CONSTRAINT ck_device_state_logs__state_valid CHECK (state IN ('online', 'offline'))`).Error,
				db.Exec(`ALTER TABLE device_state_notifications ADD CONSTRAINT ck_device_state_notifications__state_valid CHECK (state IN ('online', 'offline'))`).Error,
				db.Exec(`ALTER TABLE device_state_notifications ADD CONSTRAINT ck_device_state_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE device_limit_notifications ADD CONSTRAINT ck_device_limit_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE device_discovery_notifications ADD CONSTRAINT ck_device_discovery_notifications__processed_valid CHECK (processed IN (0, 1))`).Error,

				db.Exec(`ALTER TABLE vendor_lookups ADD CONSTRAINT ck_vendor_lookups__processed_valid CHECK (processed IN (0, 1))`).Error,
				db.Exec(`ALTER TABLE vendor_lookups ADD CONSTRAINT ck_vendor_lookups__succeeded_valid CHECK (succeeded IN (0, 1))`).Error,

				db.Exec(`ALTER TABLE users ADD CONSTRAINT ck_users__role_valid CHECK (role IN (1, 2, 3))`).Error,
				db.Exec(`ALTER TABLE users ADD CONSTRAINT ck_users__last_login_at_valid CHECK (last_login_at >= created_at)`).Error,
			)
		},
	},
}
