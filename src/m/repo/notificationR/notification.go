package notificationR

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/pager"
	"crdx.org/lighthouse/util/dbutil"
)

type ListView struct {
	CreatedAt time.Time
	Subject   string
	Body      string
}

func GetListViewRowCount() uint {
	return db.Query[uint](`
		SELECT count(*)
		FROM notifications
	`)
}

func GetListView(page uint, perPage uint) []ListView {
	q := dbutil.NewQueryBuilder(`
		SELECT
			created_at,
			subject,
			body
		FROM notifications
		ORDER BY created_at DESC
	`)

	q.Append(`LIMIT ?, ?`, pager.GetOffset(page, perPage), perPage)

	return db.Query[[]ListView](q.Query(), q.Args()...)
}
