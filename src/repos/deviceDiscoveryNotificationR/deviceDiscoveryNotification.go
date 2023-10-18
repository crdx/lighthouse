package deviceDiscoveryNotificationR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Unprocessed() []*m.DeviceDiscoveryNotification {
	return db.B[m.DeviceDiscoveryNotification]().
		Where("processed = 0").
		Order("created_at ASC").
		Find()
}
