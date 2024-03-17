package service

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
)

type discovery struct {
	Notification *db.DeviceServiceNotification
	Service      *db.Service
	Device       *db.Device
}

func Notifications() *db.Notification {
	notifications := db.FindUnprocessedDeviceServiceNotifications()
	if len(notifications) == 0 {
		return nil
	}

	var discoveries []*discovery

	for _, notification := range notifications {
		device, foundDevice := db.FindDevice(notification.DeviceID)
		service, foundService := db.FindService(notification.ServiceID.V)

		if !foundDevice || !foundService {
			notification.Delete()
			continue
		}

		discoveries = append(discoveries, &discovery{
			Notification: notification,
			Device:       device,
			Service:      service,
		})
	}

	defer func() {
		for _, notification := range notifications {
			notification.UpdateProcessed(true)
		}
	}()

	subject := getSubject(discoveries)
	body := getBody(discoveries)

	return &db.Notification{
		Subject: subject,
		Body:    body,
	}
}

func getSubject(discoveries []*discovery) string {
	if len(discoveries) == 1 {
		return fmt.Sprintf("new service found on %s", discoveries[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d new services found", len(discoveries))
	}
}

func getBody(discoveries []*discovery) string {
	var s strings.Builder
	for _, discovery := range discoveries {
		s.WriteString(fmt.Sprintf("%s\n", discovery.Service.Details()))
	}
	return strings.TrimSpace(s.String())
}
