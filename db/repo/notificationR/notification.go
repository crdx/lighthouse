package notificationR

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/pager"
)

func GetListRowCount() int {
	return int(db.CountNotificationsListView())
}

func GetList(page int, perPage int) []*db.Notification {
	return db.FindNotificationsListView(
		int64(pager.GetOffset(page, perPage)),
		int64(perPage),
	)
}
