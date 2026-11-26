package notifier

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
)

type discovery struct {
	Notification *db.DeviceDiscoveryNotification
	Device       *db.Device
}

func discoveryNotifications() (*db.Notification, func()) {
	notifications := db.FindUnprocessedDeviceDiscoveryNotifications()
	if len(notifications) == 0 {
		return nil, nil
	}

	var discoveries []*discovery

	for _, notification := range notifications {
		if device, found := db.FindDevice(notification.DeviceID); !found {
			notification.Delete()
			continue
		} else {
			discoveries = append(discoveries, &discovery{
				Notification: notification,
				Device:       device,
			})
		}
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

	subject := discoverySubject(discoveries)
	body := discoveryBody(discoveries)

	return &db.Notification{
		Subject: subject,
		Body:    body,
	}, markProcessed
}

func discoverySubject(discoveries []*discovery) string {
	if len(discoveries) == 1 {
		return fmt.Sprintf("%s joined the network", discoveries[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d devices joined the network", len(discoveries))
	}
}

func discoveryBody(discoveries []*discovery) string {
	var s strings.Builder
	for _, discovery := range discoveries {
		fmt.Fprintf(&s, "%s\n", discovery.Device.Details())
	}
	return strings.TrimSpace(s.String())
}
