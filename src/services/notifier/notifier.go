// Package notifier batches up and sends out pending notifications.
package notifier

import (
	"log/slog"

	"crdx.org/lighthouse/services"
	"crdx.org/lighthouse/services/notifier/discovery"
	"crdx.org/lighthouse/services/notifier/state"
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
	state.ProcessNotifications()
	discovery.ProcessNotifications()

	return nil
}
