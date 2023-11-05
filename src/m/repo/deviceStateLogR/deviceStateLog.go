package deviceStateLogR

import (
	"database/sql"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/pager"
	"crdx.org/lighthouse/util"
)

func All() []*m.DeviceStateLog {
	return db.B[m.DeviceStateLog]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func LatestActivityForDevice(deviceID uint) []*m.DeviceStateLog {
	return db.B[m.DeviceStateLog]().
		Where("device_id = ?", deviceID).
		Limit(6).
		Order("created_at DESC").
		Find()
}

type ListView struct {
	CreatedAt time.Time
	DeviceID  string
	Name      string
	DeletedAt sql.NullTime
	Icon      string
	State     string
}

func (self ListView) IconClass() string {
	return util.IconToClass(self.Icon)
}

func GetListViewRowCount(deviceID uint) uint {
	q := db.Q(`
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
	q := db.Q(`
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
