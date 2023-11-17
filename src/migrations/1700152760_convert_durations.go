package migrations

import (
	"errors"

	"crdx.org/db"
	"gorm.io/gorm"
)

func ConvertDurations(id string) *db.Migration {
	return &db.Migration{
		ID:   id,
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
	}
}
