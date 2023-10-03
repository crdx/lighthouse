package deviceM

import (
	"fmt"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/models/adapterM"
	"gorm.io/gorm"
)

const (
	StateOnline  = "online"
	StateOffline = "offline"
)

// https://gorm.io/docs/models.html
type Device struct {
	ID                   uint           `gorm:"primarykey"`
	CreatedAt            time.Time      `gorm:""`
	UpdatedAt            time.Time      `gorm:""`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	NetworkID            uint           `gorm:"not null"`
	Name                 string         `gorm:"size:255;not null"`
	Hostname             string         `gorm:"size:255;not null"`
	State                string         `gorm:"size:32;not null"`
	LastSeen             time.Time      `gorm:""`
	NotifyOnStatusChange bool           `gorm:"not null;default:false"`
	GracePeriod          uint           `gorm:"not null;default:300"`
}

func (self *Device) Update(values ...any) {
	For(self.ID).Update(values...)
}

func (self *Device) Fresh() *Device {
	i, _ := For(self.ID).First()
	return i
}

func (self *Device) Adapters() []*adapterM.Adapter {
	return db.B(adapterM.Adapter{DeviceID: self.ID}).Find()
}

func (self *Device) Delete() {
	for _, adapter := range self.Adapters() {
		adapter.Delete()
	}

	For(self.ID).Delete()
}

func For(id uint) *db.Builder[Device] {
	return db.B(Device{ID: id})
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

func GetListView(sortColumn string, sortDirection string) []*DeviceListView {
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
		return []*DeviceListView{}
	}

	orderBy := fmt.Sprintf(orderByTemplate, sortDirection)

	return db.Query[[]*DeviceListView](fmt.Sprintf(`
		SELECT
			D.id,
			D.name,
			D.state,
			D.hostname,
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
