package service

import (
	"fmt"
	"strings"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceServiceNotificationR"
)

type discovery struct {
	Notification *m.DeviceServiceNotification
	Service      *m.Service
	Device       *m.Device
}

func Notifications() *m.Notification {
	notifications := deviceServiceNotificationR.Unprocessed()
	if len(notifications) == 0 {
		return nil
	}

	var discoveries []*discovery

	for _, notification := range notifications {
		device, foundDevice := db.First[m.Device](notification.DeviceID)
		service, foundService := db.First[m.Service](notification.ServiceID)

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
			notification.Update("processed", true)
		}
	}()

	subject := getSubject(discoveries)
	body := getBody(discoveries)

	return &m.Notification{
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
