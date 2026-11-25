// Package notifier batches up and sends out pending notifications.
package notifier

import (
	"log/slog"

	"crdx.org/lighthouse/cmd/lighthouse/services"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier/discovery"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier/limit"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier/service"
	"crdx.org/lighthouse/cmd/lighthouse/services/notifier/state"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/util/mailutil"
)

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
	self.add(discovery.Notifications())
	self.add(limit.Notifications())
	self.add(service.Notifications())
	self.add(state.Notifications())

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
