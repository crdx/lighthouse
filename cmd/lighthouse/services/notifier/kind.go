package notifier

import "crdx.org/lighthouse/db"

// row is a stored notification of one kind, such as a device state notification.
type row interface {
	Delete() bool
	UpdateProcessed(value bool) bool
}

// source gives the pending notification of one kind.
type source interface {
	pending() (*db.Notification, func())
}

// notificationKind describes one kind of notification. It tells how to load the rows of the kind,
// and how to word one mail about them.
type notificationKind[Row row, Item any] struct {
	// find returns the unprocessed rows, oldest first.
	find func() []Row
	// resolve returns the item for a row, or false if the device or service of the row is gone.
	resolve func(Row) (Item, bool)
	// keep returns the items to send. A nil keep sends every item.
	keep    func(items []Item) []Item
	subject func(kept []Item) string
	body    func(kept []Item, all []Item) string
}

// pending returns the notification for the unprocessed rows, and a function that marks the rows as
// processed. It deletes a row whose device or service is gone. It marks the rows itself when there
// is nothing to send.
func (self notificationKind[Row, Item]) pending() (*db.Notification, func()) {
	rows := self.find()
	if len(rows) == 0 {
		return nil, nil
	}

	var items []Item
	for _, row := range rows {
		if item, found := self.resolve(row); found {
			items = append(items, item)
		} else {
			row.Delete()
		}
	}

	markProcessed := func() {
		for _, row := range rows {
			row.UpdateProcessed(true)
		}
	}

	kept := items
	if self.keep != nil {
		kept = self.keep(items)
	}

	if len(kept) == 0 {
		markProcessed()
		return nil, nil
	}

	return &db.Notification{
		Subject: self.subject(kept),
		Body:    self.body(kept, items),
	}, markProcessed
}
