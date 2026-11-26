package notifier

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
)

// discoveredDevices notifies about devices that joined the network.
var discoveredDevices = notificationKind[*db.DeviceDiscoveryNotification, *db.Device]{
	find: db.FindUnprocessedDeviceDiscoveryNotifications,
	resolve: func(notification *db.DeviceDiscoveryNotification) (*db.Device, bool) {
		return db.FindDevice(notification.DeviceID)
	},
	subject: func(devices []*db.Device) string {
		if len(devices) == 1 {
			return fmt.Sprintf("%s joined the network", devices[0].DisplayName())
		}
		return fmt.Sprintf("%d devices joined the network", len(devices))
	},
	body: func(devices []*db.Device, _ []*db.Device) string {
		var s strings.Builder
		for _, device := range devices {
			fmt.Fprintf(&s, "%s\n", device.Details())
		}
		return strings.TrimSpace(s.String())
	},
}
