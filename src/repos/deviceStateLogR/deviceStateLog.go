package deviceStateLogR

import (
	"database/sql"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/pager"
	"crdx.org/lighthouse/util/dbutil"
)

func All() []*m.DeviceStateLog {
	return db.B[m.DeviceStateLog]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func LatestActivityForDevice(deviceID uint) []*m.DeviceStateLog {
	return db.B(m.DeviceStateLog{DeviceID: deviceID}).Limit(5).Order("created_at DESC").Find()
}

type ListView struct {
	CreatedAt time.Time
	DeviceID  string
	Name      string
	DeletedAt sql.NullTime
	Icon      string
	State     string
}

func GetListViewTotal(deviceID uint) uint {
	q := dbutil.NewQueryBuilder(`
		SELECT count(*)
		FROM device_state_logs DSL
		INNER JOIN devices D on D.id = DSL.device_id
	`)

	if deviceID != 0 {
		q.And("device_id = ?", deviceID)
	}

	return db.Query[uint](q.Query(), q.Args()...)
}

func GetListView(deviceID uint, page uint, perPage uint) []ListView {
	q := dbutil.NewQueryBuilder(`
		SELECT
			DSL.created_at,
			DSL.device_id,
			D.name,
			D.icon,
			D.deleted_at,
			DSL.state
		FROM device_state_logs DSL
		INNER JOIN devices D on D.id = DSL.device_id
	`)

	if deviceID != 0 {
		q.And("device_id = ?", deviceID)
	}

	q.Append("ORDER BY created_at DESC")
	q.Append(`LIMIT ?, ?`, pager.GetOffset(page, perPage), perPage)

	return db.Query[[]ListView](q.Query(), q.Args()...)
}
