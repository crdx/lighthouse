package conf

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrations = []*gormigrate.Migration{
	{
		ID: "1697905074_InsertDefaultTimeZone",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(`
				INSERT INTO settings
					(name, value, created_at, updated_at)
				VALUES
					('timezone', 'Europe/London', UTC_TIMESTAMP(), UTC_TIMESTAMP())
			`).Error
		},
	},
}
