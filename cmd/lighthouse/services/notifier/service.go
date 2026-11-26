package notifier

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
)

type newService struct {
	Notification *db.DeviceServiceNotification
	Service      *db.Service
	Device       *db.Device
}

func serviceNotifications() (*db.Notification, func()) {
	notifications := db.FindUnprocessedDeviceServiceNotifications()
	if len(notifications) == 0 {
		return nil, nil
	}

	var discoveries []*newService

	for _, notification := range notifications {
		device, foundDevice := db.FindDevice(notification.DeviceID)
		service, foundService := db.FindService(notification.ServiceID.V)

		if !foundDevice || !foundService {
			notification.Delete()
			continue
		}

		discoveries = append(discoveries, &newService{
			Notification: notification,
			Device:       device,
			Service:      service,
		})
	}

	markProcessed := func() {
		for _, notification := range notifications {
			notification.UpdateProcessed(true)
		}
	}

	if len(discoveries) == 0 {
		markProcessed()
		return nil, nil
	}

	subject := serviceSubject(discoveries)
	body := serviceBody(discoveries)

	return &db.Notification{
		Subject: subject,
		Body:    body,
	}, markProcessed
}

func serviceSubject(discoveries []*newService) string {
	if len(discoveries) == 1 {
		return fmt.Sprintf("new service found on %s", discoveries[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d new services found", len(discoveries))
	}
}

func serviceBody(discoveries []*newService) string {
	var s strings.Builder
	for _, discovery := range discoveries {
		fmt.Fprintf(&s, "%s\n", discovery.Service.Details())
	}
	return strings.TrimSpace(s.String())
}
