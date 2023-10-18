package state

import (
	"fmt"
	"strings"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/repos/deviceStateNotificationR"
	"crdx.org/lighthouse/util/mailutil"
	"crdx.org/lighthouse/util/timeutil"
	"github.com/samber/lo"
)

type DeviceID = uint

type transition struct {
	Notification *m.DeviceStateNotification
	Device       *m.Device
}

func (self *transition) String() string {
	if self.Notification.State == deviceR.StateOnline {
		return fmt.Sprintf("%s is online", self.Device.Identifier())
	} else if self.Notification.State == deviceR.StateOffline {
		return fmt.Sprintf("%s is offline", self.Device.Identifier())
	} else {
		return fmt.Sprintf(
			"%s transitioned to an unknown state (%s)",
			self.Device.Identifier(),
			self.Notification.State,
		)
	}
}

func (self *transition) TimestampedString() string {
	return fmt.Sprintf("%s — %s", timeutil.ToLocal(self.Notification.CreatedAt).Format("15:04"), self.String())
}

func ProcessNotifications() {
	notifications := deviceStateNotificationR.Unprocessed()
	if len(notifications) == 0 {
		return
	}

	var transitions []*transition

	for _, notification := range notifications {
		if device, found := db.First[m.Device](notification.DeviceID); !found {
			notification.Delete()
			continue
		} else {
			transitions = append(transitions, &transition{
				Notification: notification,
				Device:       device,
			})
		}
	}

	defer func() {
		for _, transition := range transitions {
			transition.Notification.Update("processed", true)
		}
	}()

	if len(transitions) == 0 {
		return
	}

	newTransitions := getNewTransitions(transitions)

	if len(newTransitions) == 0 {
		return
	}

	subject := getSubject(newTransitions)
	body := getBody(newTransitions, transitions)

	// Panic here as this will probably be a recoverable failure e.g. intermittent network failure.
	lo.Must0(mailutil.SendNotification(subject, body))
}

func getNewTransitions(transitions []*transition) []*transition {
	initialTransitions := map[DeviceID]*transition{}
	finalTransitions := map[DeviceID]*transition{}

	// Notifications are looped through in chronological order, so each iteration of the loop will
	// overwrite the map entry in finalTransitions for the same device — this is intentional. This
	// means the final entry will be the most recent notification. We also store the first one in
	// initialTransitions, as this will let us compare the first and last notification below.
	for _, transition := range transitions {
		deviceID := transition.Device.ID

		if _, ok := initialTransitions[deviceID]; !ok {
			initialTransitions[deviceID] = transition
		}

		finalTransitions[deviceID] = transition
	}

	// If the final transition is the same as the first one then we know the device has overall
	// changed state. This logic holds if there is only one transition logged for this period, as
	// then it's being compared to itself.
	return lo.Values(lo.OmitBy(finalTransitions, func(deviceID DeviceID, finalTransition *transition) bool {
		return finalTransition.Notification.State != initialTransitions[deviceID].Notification.State
	}))
}

func getSubject(newTransitions []*transition) string {
	if len(newTransitions) == 0 {
		return "device state has fluctuated"
	}

	if len(newTransitions) == 1 {
		return newTransitions[0].TimestampedString()
	}

	isOnline := func(t *transition) bool { return t.Notification.State == deviceR.StateOnline }
	isOffline := func(t *transition) bool { return t.Notification.State == deviceR.StateOffline }

	if lo.EveryBy(newTransitions, isOnline) {
		return fmt.Sprintf("%d devices are online", len(newTransitions))
	}

	if lo.EveryBy(newTransitions, isOffline) {
		return fmt.Sprintf("%d devices are offline", len(newTransitions))
	}

	return fmt.Sprintf(
		"%d devices changed state (%d online, %d offline)",
		len(newTransitions),
		lo.CountBy(newTransitions, isOnline),
		lo.CountBy(newTransitions, isOffline),
	)
}

func getBody(newTransitions []*transition, allTransitions []*transition) string {
	var s strings.Builder

	if len(newTransitions) > 1 {
		for _, transition := range newTransitions {
			s.WriteString(transition.TimestampedString() + "\n")
		}
	}

	if len(allTransitions) > 0 && len(allTransitions) > len(newTransitions) {
		if len(newTransitions) > 0 {
			s.WriteString("\nAll activity:\n\n")
		}

		for _, transition := range allTransitions {
			s.WriteString(transition.TimestampedString() + "\n")
		}
	}

	return s.String()
}
