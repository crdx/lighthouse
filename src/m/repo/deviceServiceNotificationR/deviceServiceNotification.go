package deviceServiceNotificationR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Unprocessed() []*m.DeviceServiceNotification {
	return db.B[m.DeviceServiceNotification]().
		Where("processed = 0").
		Order("created_at ASC").
		Find()
}
