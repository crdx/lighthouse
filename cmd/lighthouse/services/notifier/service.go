package notifier

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
)

type newService struct {
	Service *db.Service
	Device  *db.Device
}

// newServices notifies about services found on devices.
var newServices = notificationKind[*db.DeviceServiceNotification, *newService]{
	find: db.FindUnprocessedDeviceServiceNotifications,
	resolve: func(notification *db.DeviceServiceNotification) (*newService, bool) {
		device, foundDevice := db.FindDevice(notification.DeviceID)
		service, foundService := db.FindService(notification.ServiceID.V)
		return &newService{Service: service, Device: device}, foundDevice && foundService
	},
	subject: func(services []*newService) string {
		if len(services) == 1 {
			return fmt.Sprintf("new service found on %s", services[0].Device.DisplayName())
		}
		return fmt.Sprintf("%d new services found", len(services))
	},
	body: func(services []*newService, _ []*newService) string {
		var s strings.Builder
		for _, service := range services {
			fmt.Fprintf(&s, "%s\n", service.Service.Details())
		}
		return strings.TrimSpace(s.String())
	},
}
