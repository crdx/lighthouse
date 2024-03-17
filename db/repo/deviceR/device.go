package deviceR

import (
	"fmt"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/util"
)

const (
	StateOnline  = "online"
	StateOffline = "offline"
)

type ListItem struct {
	ID         int64
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

func (self *ListItem) IconClass() string {
	return util.IconToClass(self.Icon)
}

func GetList(sortColumn string, sortDirection string, filter string) []ListItem {
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
		"online":     fmt.Sprintf("AND state = '%s'", StateOnline),
		"offline":    fmt.Sprintf("AND state = '%s'", StateOffline),
		"watched":    "AND watch = 1",
		"unwatched":  "AND watch = 0",
		"pingable":   "AND ping = 1",
		"unpingable": "AND ping = 0",
		"all":        "",
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
	rows, err := db.Query(fmt.Sprintf(`
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
	if err != nil {
		return nil
	}

	var list []ListItem
	for rows.Next() {
		var item ListItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.State,
			&item.Hostname,
			&item.Icon,
			&item.Watch,
			&item.MACAddress,
			&item.Vendor,
			&item.IPAddress,
			&item.LastSeenAt,
		); err != nil {
			return nil
		}
		list = append(list, item)
	}

	return list
}

type Counters struct {
	All        int
	Online     int
	Offline    int
	Watched    int
	Unwatched  int
	Pingable   int
	Unpingable int
}

func GetListCounts() *Counters {
	return &Counters{
		All:        len(db.FindDevices()),
		Online:     len(db.FindDevicesByState(StateOnline)),
		Offline:    len(db.FindDevicesByState(StateOffline)),
		Watched:    len(db.FindDevicesByWatch(true)),
		Unwatched:  len(db.FindDevicesByWatch(false)),
		Pingable:   len(db.FindDevicesByPing(true)),
		Unpingable: len(db.FindDevicesByPing(false)),
	}
}
