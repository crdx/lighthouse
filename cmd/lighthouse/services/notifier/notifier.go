// Package notifier batches up and sends out pending notifications.
package notifier

import (
	"log/slog"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/util/mailutil"
)

// sources are the kinds of notification, in the order the notifier sends them.
var sources = []source{discoveredDevices, devicesOverLimit, newServices, stateChanges}

type Notifier struct {
	logger *slog.Logger
}

func New() *Notifier {
	return &Notifier{}
}

func (self *Notifier) Init(args *services.Args) error {
	self.logger = args.Logger
	return nil
}

func (self *Notifier) Run() error {
	for _, source := range sources {
		self.add(source.pending())
	}

	return nil
}

func (self *Notifier) add(notification *db.Notification, markProcessed func()) {
	if notification == nil {
		return
	}

	if err := mailutil.Send(notification.Subject, notification.Body); err != nil {
		self.logger.Error("unable to send notification, will retry", "subject", notification.Subject, "err", err)
		return
	}

	db.CreateNotification(notification)
	markProcessed()
}
