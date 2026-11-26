package notifier

import (
	"fmt"
	"strings"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/util/timeutil"
)

type trespasser struct {
	Notification *db.DeviceLimitNotification
	Device       *db.Device
}

func (self *trespasser) String() string {
	return fmt.Sprintf(
		"%s is still online since %s (%s)",
		self.Device.Identifier(),
		timeutil.ToLocal(self.Notification.StateUpdatedAt).Format(constants.TimeFormatReadable),
		timeutil.FormatDuration(time.Since(self.Notification.StateUpdatedAt), true, 1, "ago"),
	)
}

func limitNotifications() (*db.Notification, func()) {
	notifications := db.FindUnprocessedDeviceLimitNotifications()
	if len(notifications) == 0 {
		return nil, nil
	}

	var trespassers []*trespasser

	for _, notification := range notifications {
		if device, found := db.FindDevice(notification.DeviceID); !found {
			notification.Delete()
			continue
		} else {
			trespassers = append(trespassers, &trespasser{
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

	if len(trespassers) == 0 {
		markProcessed()
		return nil, nil
	}

	subject := limitSubject(trespassers)
	body := limitBody(trespassers)

	return &db.Notification{
		Subject: subject,
		Body:    body,
	}, markProcessed
}

func limitSubject(trespassers []*trespasser) string {
	if len(trespassers) == 1 {
		return fmt.Sprintf("%s is still online", trespassers[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d devices are still online", len(trespassers))
	}
}

func limitBody(trespassers []*trespasser) string {
	var s strings.Builder

	for _, trespasser := range trespassers {
		s.WriteString(trespasser.String())
		s.WriteString("\n")
	}

	return strings.TrimSpace(s.String())
}
