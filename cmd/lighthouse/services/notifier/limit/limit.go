package limit

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

func Notifications() *db.Notification {
	notifications := db.FindUnprocessedDeviceLimitNotifications()
	if len(notifications) == 0 {
		return nil
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

	defer func() {
		for _, notification := range notifications {
			notification.UpdateProcessed(true)
		}
	}()

	subject := getSubject(trespassers)
	body := getBody(trespassers)

	return &db.Notification{
		Subject: subject,
		Body:    body,
	}
}

func getSubject(trespassers []*trespasser) string {
	if len(trespassers) == 1 {
		return fmt.Sprintf("%s is still online", trespassers[0].Device.DisplayName())
	} else {
		return fmt.Sprintf("%d devices are still online", len(trespassers))
	}
}

func getBody(trespassers []*trespasser) string {
	var s strings.Builder

	for _, trespasser := range trespassers {
		s.WriteString(trespasser.String() + "\n")
	}

	return strings.TrimSpace(s.String())
}
