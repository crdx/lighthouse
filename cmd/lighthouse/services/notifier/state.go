package notifier

import (
	"fmt"
	"strings"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceR"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/util/timeutil"
	"github.com/samber/lo"
)

type transition struct {
	Notification *db.DeviceStateNotification
	Device       *db.Device
}

func (self *transition) String() string {
	return fmt.Sprintf("%s is %s", self.Device.Identifier(), self.Notification.State)
}

func (self *transition) Verb() string {
	return map[string]string{
		deviceR.StateOnline:  "came",
		deviceR.StateOffline: "went",
	}[self.Notification.State]
}

func (self *transition) TimestampedString() string {
	return fmt.Sprintf(
		"%s %s %s at %s",
		self.Device.Identifier(),
		self.Verb(),
		self.Notification.State,
		timeutil.ToLocal(self.Notification.CreatedAt).Format(constants.TimeFormatEuropeanKitchen),
	)
}

// stateChanges notifies about devices that went online or offline.
var stateChanges = notificationKind[*db.DeviceStateNotification, *transition]{
	find: db.FindUnprocessedDeviceStateNotifications,
	resolve: func(notification *db.DeviceStateNotification) (*transition, bool) {
		device, found := db.FindDevice(notification.DeviceID)
		return &transition{Notification: notification, Device: device}, found
	},
	keep:    changedDevices,
	subject: stateSubject,
	body:    stateBody,
}

// changedDevices returns the last transition of each device that ends in a different state from
// the one it started in. A device that went offline and came back is left out.
func changedDevices(all []*transition) []*transition {
	initialTransitions := map[int64]*transition{}
	finalTransitions := map[int64]*transition{}

	// The transitions are in chronological order, so the last one for a device overwrites the
	// others in finalTransitions. The first one stays in initialTransitions.
	for _, transition := range all {
		deviceID := transition.Device.ID

		if _, ok := initialTransitions[deviceID]; !ok {
			initialTransitions[deviceID] = transition
		}

		finalTransitions[deviceID] = transition
	}

	// A device with one transition compares it with itself, so it is kept.
	return lo.Values(lo.OmitBy(finalTransitions, func(deviceID int64, finalTransition *transition) bool {
		return finalTransition.Notification.State != initialTransitions[deviceID].Notification.State
	}))
}

func stateSubject(changes []*transition) string {
	if len(changes) == 1 {
		return changes[0].String()
	}

	isOnline := func(change *transition) bool { return change.Notification.State == deviceR.StateOnline }
	isOffline := func(change *transition) bool { return change.Notification.State == deviceR.StateOffline }

	if lo.EveryBy(changes, isOnline) {
		return fmt.Sprintf("%d devices are online", len(changes))
	}

	if lo.EveryBy(changes, isOffline) {
		return fmt.Sprintf("%d devices are offline", len(changes))
	}

	return fmt.Sprintf(
		"%d devices changed state (%d online, %d offline)",
		len(changes),
		lo.CountBy(changes, isOnline),
		lo.CountBy(changes, isOffline),
	)
}

// stateBody lists the changes. If some devices only went offline and came back, it also lists all
// the activity.
func stateBody(changes []*transition, all []*transition) string {
	var s strings.Builder

	for _, change := range changes {
		s.WriteString(change.TimestampedString())
		s.WriteString("\n")
	}

	if len(all) > len(changes) {
		s.WriteString("\nAll activity:\n\n")

		for _, transition := range all {
			s.WriteString(transition.TimestampedString())
			s.WriteString("\n")
		}
	}

	return s.String()
}
