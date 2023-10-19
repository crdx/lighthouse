package discovery

import (
	"fmt"
	"strings"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceDiscoveryNotificationR"
	"crdx.org/lighthouse/util/mailutil"
	"github.com/samber/lo"
)

type discovery struct {
	Notification *m.DeviceDiscoveryNotification
	Device       *m.Device
}

func ProcessNotifications() {
	notifications := deviceDiscoveryNotificationR.Unprocessed()
	if len(notifications) == 0 {
		return
	}

	var discoveries []*discovery

	for _, notification := range notifications {
		if device, found := db.First[m.Device](notification.DeviceID); !found {
			notification.Delete()
			continue
		} else {
			discoveries = append(discoveries, &discovery{
				Notification: notification,
				Device:       device,
			})
		}
	}

	defer func() {
		for _, notification := range notifications {
			notification.Update("processed", true)
		}
	}()

	subject := getSubject(discoveries)
	body := getBody(discoveries)

	// Panic here as this will probably be a recoverable failure e.g. intermittent network failure.
	lo.Must0(mailutil.Send(subject, body))
}

func getSubject(discoveries []*discovery) string {
	if len(discoveries) == 1 {
		return fmt.Sprintf("%s joined the network", discoveries[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d devices joined the network", len(discoveries))
	}
}

func getBody(discoveries []*discovery) string {
	var s strings.Builder
	for _, discovery := range discoveries {
		s.WriteString(fmt.Sprintf("%s\n", discovery.Device.Details()))
	}
	return s.String()
}
