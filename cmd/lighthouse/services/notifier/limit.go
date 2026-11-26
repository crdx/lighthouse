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

// devicesOverLimit notifies about devices that stayed online past their limit.
var devicesOverLimit = notificationKind[*db.DeviceLimitNotification, *trespasser]{
	find: db.FindUnprocessedDeviceLimitNotifications,
	resolve: func(notification *db.DeviceLimitNotification) (*trespasser, bool) {
		device, found := db.FindDevice(notification.DeviceID)
		return &trespasser{Notification: notification, Device: device}, found
	},
	subject: func(trespassers []*trespasser) string {
		if len(trespassers) == 1 {
			return fmt.Sprintf("%s is still online", trespassers[0].Device.DisplayName())
		}
		return fmt.Sprintf("%d devices are still online", len(trespassers))
	},
	body: func(trespassers []*trespasser, _ []*trespasser) string {
		var s strings.Builder
		for _, trespasser := range trespassers {
			s.WriteString(trespasser.String())
			s.WriteString("\n")
		}
		return strings.TrimSpace(s.String())
	},
}
