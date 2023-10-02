package deviceModel

import (
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/util"
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html
type Device struct {
	ID                   uint           `gorm:"primarykey"`
	CreatedAt            time.Time      `gorm:""`
	UpdatedAt            time.Time      `gorm:""`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	MACAddress           string         `gorm:"size:17;not null"`
	NetworkID            uint           `gorm:"not null"`
	Name                 string         `gorm:"size:255;not null"`
	Hostname             string         `gorm:"size:255;not null"`
	Vendor               string         `gorm:"size:255;not null"`
	State                string         `gorm:"size:32;not null"`
	LastSeen             time.Time      `gorm:""`
	NotifyOnStatusChange bool           `gorm:"not null;default:false"`
	GracePeriod          uint           `gorm:"not null;default:300"`
}

type DeviceListView struct {
	ID         uint
	Name       string
	State      string
	Hostname   string
	MACAddress string
	Vendor     string
	IPAddress  string
	LastSeen   time.Time
}

const (
	StateOnline  = "online"
	StateOffline = "offline"
)

func GetListView(sortColumn string, sortDirection string) []DeviceListView {
	// Ensure sort is stable by appending "D.id ASC" to some of these.
	orderByTemplates := map[string]string{
		"name":   "D.name %s, D.id ASC",
		"ip":     "INET_ATON(DM1.ip_address) %s",
		"vendor": "D.vendor %s, D.id ASC",
		"mac":    "D.mac_address %s",
		"seen":   "DM1.last_seen %s, D.id ASC",
	}

	orderByTemplate, ok := orderByTemplates[sortColumn]
	if !ok {
		return []DeviceListView{}
	}

	orderBy:= fmt.Sprintf(orderByTemplate, sortDirection)

	return db.Query[[]DeviceListView](fmt.Sprintf(`
		SELECT
			D.id,
			D.name,
			D.state,
			D.hostname,
			D.mac_address,
			D.vendor,
			DM1.ip_address,
			DM1.last_seen
		FROM devices D
		INNER JOIN device_mappings DM1 ON D.id = DM1.device_id
		LEFT JOIN device_mappings DM2 ON DM1.device_id = DM2.device_id AND DM1.last_seen < DM2.last_seen
		WHERE D.deleted_at IS NULL
		AND DM1.deleted_at IS NULL
		AND DM2.id IS NULL
		ORDER BY %s
	`, orderBy))
}

func Upsert(networkID uint, macAddress string) (Device, bool) {
	// This method has the potential to be called very often, so let's not hang onto a model object
	// for longer than necessary. Immediately create the record if it doesn't exist, and then just
	// run the query that updates the specific fields.
	device, found := db.FirstOrCreate(Device{
		NetworkID:  networkID,
		MACAddress: macAddress,
	})

	columns := db.Map{}

	if !found {
		columns["state"] = StateOnline
	}

	if device.Vendor == "" {
		if vendor, found := util.GetVendor(macAddress); found {
			columns["vendor"] = vendor
		}
	}

	columns["last_seen"] = time.Now()

	q := Device{ID: device.ID}

	db.B(q).Update(columns)
	device, _ = db.B(q).First()

	return device, found
}

func UpdateHostname(macAddress string, hostname string) {
	q := Device{MACAddress: macAddress}

	if device, found := db.B(q).First(); found {
		updates := map[string]any{}

		if device.Name == "" {
			updates["name"] = hostname
		}

		updates["hostname"] = hostname
		db.B(q).Update(updates)
	}
}
