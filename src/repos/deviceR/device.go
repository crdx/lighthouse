package deviceR

import (
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func All() []*m.Device {
	return db.B[m.Device]().Order("name ASC").Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

const (
	StateOnline  = "online"
	StateOffline = "offline"
)

type ListView struct {
	ID         uint
	Name       string
	State      string
	Hostname   string
	Icon       string
	MACAddress string
	Vendor     string
	IPAddress  string
	LastSeen   time.Time
}

func GetListView(sortColumn string, sortDirection string) []*ListView {
	// Ensure sort is stable by appending "D.id ASC" to some of these.
	orderByTemplates := map[string]string{
		"name":   "D.name %s, D.id ASC",
		"ip":     "INET_ATON(A1.ip_address) %s",
		"vendor": "A1.vendor %s, D.id ASC",
		"mac":    "A1.mac_address %s",
		"seen":   "A1.last_seen %s, D.id ASC",
	}

	orderByTemplate, ok := orderByTemplates[sortColumn]
	if !ok {
		return []*ListView{}
	}

	orderBy := fmt.Sprintf(orderByTemplate, sortDirection)

	return db.Query[[]*ListView](fmt.Sprintf(`
		SELECT
			D.id,
			D.name,
			D.state,
			D.hostname,
			D.icon,
			A1.mac_address,
			A1.vendor,
			A1.ip_address,
			A1.last_seen
		FROM devices D
		INNER JOIN adapters A1 ON D.id = A1.device_id
		LEFT JOIN adapters A2 ON A1.device_id = A2.device_id AND A1.last_seen < A2.last_seen
		WHERE D.deleted_at IS NULL
		AND A1.deleted_at IS NULL
		AND A2.id IS NULL
		ORDER BY %s
	`, orderBy))
}
