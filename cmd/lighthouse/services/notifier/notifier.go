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
	"github.com/samber/lo"
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

func (*Notifier) Run() error {
	add(discovery.Notifications())
	add(limit.Notifications())
	add(service.Notifications())
	add(state.Notifications())

	return nil
}

func add(notification *db.Notification) {
	if notification == nil {
		return
	}

	db.CreateNotification(notification)
	lo.Must0(mailutil.Send(notification.Subject, notification.Body))
}
