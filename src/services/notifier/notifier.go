// Package notifier batches up and sends out pending notifications.
package notifier

import (
	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/util/mailutil"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier/discovery"
	"crdx.org/lighthouse/services/notifier/limit"
	"crdx.org/lighthouse/services/notifier/service"
	"crdx.org/lighthouse/services/notifier/state"
	"github.com/samber/lo"
)

type Notifier struct {
	log *slog.Logger
}

func New() *Notifier {
	return &Notifier{}
}

func (self *Notifier) Init(args *services.Args) error {
	self.log = args.Logger
	return nil
}

func (*Notifier) Run() error {
	add(discovery.Notifications())
	add(limit.Notifications())
	add(service.Notifications())
	add(state.Notifications())

	return nil
}

func add(notification *m.Notification) {
	if notification == nil {
		return
	}

	db.Save(&notification)
	lo.Must0(mailutil.Send(notification.Subject, notification.Body))
}
