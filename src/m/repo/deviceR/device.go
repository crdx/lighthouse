package deviceR

import (
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/util"
)

const (
	StateOnline  = "online"
	StateOffline = "offline"
)

func All() []*m.Device {
	return db.B[m.Device]().Order("name ASC").Find()
}

type ListView struct {
	ID         uint
	Name       string
	State      string
	Hostname   string
	Icon       string
	Watch      bool
	MACAddress string
	Vendor     string
	IPAddress  string
	LastSeenAt time.Time
}

func (self ListView) IconClass() string {
	return util.IconToClass(self.Icon)
}

func GetListView(sortColumn string, sortDirection string, filter string) []ListView {
	// Ensure sort is stable by appending "D.id ASC" to some of these.
	orderByTemplates := map[string]string{
		"name":   "D.name %s, D.id ASC",
		"ip":     "INET_ATON(A1.ip_address) %s",
		"vendor": "A1.vendor %s, D.id ASC",
		"mac":    "A1.mac_address %s",
		"seen":   "A1.last_seen_at %s, D.id ASC",
		"watch":  "D.watch %s, D.id ASC",
		"type":   "D.icon %s, D.id ASC",
	}

	filters := map[string]string{
		"online":    "AND state = 'online'",
		"offline":   "AND state = 'offline'",
		"watched":   "AND watch = 1",
		"unwatched": "AND watch = 0",
		"all":       "",
	}

	orderByTemplate, orderOK := orderByTemplates[sortColumn]
	filterBy, filterOK := filters[filter]

	if !orderOK || !filterOK {
		return nil
	}

	orderBy := fmt.Sprintf(orderByTemplate, sortDirection)

	// The left join with adapters on last_seen_at finds the adapter with the newest last_seen_at
	// date. This works because the row where there is no newer last_seen_at date will contain
	// nulls for the A2 part of the table, and the where clause requires A2.id to be null.
	return db.Query[[]ListView](fmt.Sprintf(`
		SELECT
			D.id,
			D.name,
			D.state,
			D.hostname,
			D.icon,
			D.watch,
			A1.mac_address,
			A1.vendor,
			A1.ip_address,
			A1.last_seen_at
		FROM devices D
		INNER JOIN adapters A1 ON D.id = A1.device_id
		LEFT JOIN adapters A2 ON A1.device_id = A2.device_id AND A1.last_seen_at < A2.last_seen_at
		WHERE D.deleted_at IS NULL
		AND A1.deleted_at IS NULL
		AND A2.id IS NULL
		%s
		ORDER BY %s
	`, filterBy, orderBy))
}

type Counts struct {
	All       uint
	Online    uint
	Offline   uint
	Watched   uint
	Unwatched uint
}

func GetCounts() *Counts {
	return &Counts{
		All:       uint(db.B[m.Device]().Count()),
		Online:    uint(db.B[m.Device]("state = ?", "online").Count()),
		Offline:   uint(db.B[m.Device]("state = ?", "offline").Count()),
		Watched:   uint(db.B[m.Device]("watch = 1").Count()),
		Unwatched: uint(db.B[m.Device]("watch = 0").Count()),
	}
}
