package notificationR

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/pager"
)

type ListView struct {
	CreatedAt time.Time
	Subject   string
	Body      string
}

func GetListViewRowCount() int {
	return db.Query[int](`
		SELECT count(*)
		FROM notifications
		WHERE deleted_at IS NULL
	`)
}

func GetListView(page int, perPage int) []ListView {
	q := db.Q(`
		SELECT
			created_at,
			subject,
			body
		FROM notifications
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)

	q.Append(`LIMIT ?, ?`, pager.GetOffset(page, perPage), perPage)

	return db.Query[[]ListView](q.Query(), q.Args()...)
}
