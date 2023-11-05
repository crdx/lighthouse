package limit

import (
	"fmt"
	"strings"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/deviceLimitNotificationR"
	"crdx.org/lighthouse/util/timeutil"
)

type trespasser struct {
	Notification *m.DeviceLimitNotification
	Device       *m.Device
}

func (self *trespasser) String() string {
	return fmt.Sprintf(
		"%s is still online since %s (%s)",
		self.Device.Identifier(),
		timeutil.ToLocal(self.Notification.StateUpdatedAt).Format(constants.TimeFormatReadable),
		timeutil.TimeAgo(int(time.Now().Sub(self.Notification.StateUpdatedAt).Seconds()), true, 1),
	)
}

func Notifications() *m.Notification {
	notifications := deviceLimitNotificationR.Unprocessed()
	if len(notifications) == 0 {
		return nil
	}

	var trespassers []*trespasser

	for _, notification := range notifications {
		if device, found := db.First[m.Device](notification.DeviceID); !found {
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
			notification.Update("processed", true)
		}
	}()

	subject := getSubject(trespassers)
	body := getBody(trespassers)

	return &m.Notification{
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
