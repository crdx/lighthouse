package deviceLimitNotificationR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Unprocessed() []*m.DeviceLimitNotification {
	return db.B[m.DeviceLimitNotification]().
		Where("processed = 0").
		Order("created_at ASC").
		Find()
}
