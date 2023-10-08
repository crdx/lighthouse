package conf

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrations = []*gormigrate.Migration{
	{
		ID: "1696599187_UpdateGracePeriod",
		Migrate: func(db *gorm.DB) error {
			return db.Exec("UPDATE devices SET grace_period = 5").Error
		},
	},
}
