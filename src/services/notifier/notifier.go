// Package notifier batches up and sends out pending notifications.
package notifier

import (
	"log/slog"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier/discovery"
	"crdx.org/lighthouse/services/notifier/state"
	"crdx.org/lighthouse/util/mailutil"
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
	add(state.Notifications())
	add(discovery.Notifications())

	return nil
}

func add(notification *m.Notification) {
	if notification == nil {
		return
	}

	db.Save(&notification)

	if !settingR.GetBool(settingR.EnableNotifications) {
		return
	}

	lo.Must0(mailutil.Send(notification.Subject, notification.Body))
}
