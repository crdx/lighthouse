package deviceStateNotificationR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Unprocessed() []*m.DeviceStateNotification {
	return db.B[m.DeviceStateNotification]().
		Where("processed = 0").
		Order("created_at ASC").
		Find()
}
