package notificationR

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/pkg/pager"
)

type List struct {
	CreatedAt time.Time
	Subject   string
	Body      string
}

func GetListRowCount() int {
	return db.Query[int](`
		SELECT count(*)
		FROM notifications
		WHERE deleted_at IS NULL
	`)
}

func GetList(page int, perPage int) []List {
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

	return db.Query[[]List](q.Query(), q.Args()...)
}
