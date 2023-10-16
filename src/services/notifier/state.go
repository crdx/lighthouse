package notifier

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/util/timeutil"
)

type DeviceID = uint

type transition struct {
	Notification *m.DeviceStateNotification
	Device       *m.Device
}

func (self *transition) String() string {
	if self.Notification.State == deviceR.StateOnline {
		return fmt.Sprintf("%s is online", self.Device.DisplayName())
	} else if self.Notification.State == deviceR.StateOffline {
		return fmt.Sprintf("%s is offline", self.Device.DisplayName())
	} else {
		return fmt.Sprintf(
			"%s transitioned to an unknown state (%s)",
			self.Device.DisplayName(),
			self.Notification.State,
		)
	}
}

func (self *transition) TimestampedString() string {
	return fmt.Sprintf("[%s] %s", timeutil.ToLocal(self.Notification.CreatedAt).Format("15:04"), self.String())
}
