package adapterR

import (
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/util/netutil"
)

func All() []*m.Adapter {
	return db.B[m.Adapter]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func AllWithoutVendor() []*m.Adapter {
	return db.B[m.Adapter]().Where(`vendor = ""`).Find()
}

func FindByMACAddress(macAddress string) (*m.Adapter, bool) {
	return db.B(m.Adapter{MACAddress: macAddress}).First()
}

func Upsert(macAddress string, ipAddress string) (*m.Adapter, bool) {
	adapter, adapterFound := db.FirstOrCreate(m.Adapter{
		MACAddress: macAddress,
	})

	columns := db.Map{
		"last_seen":  time.Now(),
		"ip_address": ipAddress,
	}

	if adapter.Vendor == "" {
		if vendor, vendorFound := netutil.GetVendor(adapter.MACAddress); vendorFound {
			columns["vendor"] = vendor

			if !adapterFound {
				columns["name"] = vendor
			}
		}
	}

	adapter.Update(columns)

	return adapter.Fresh(), adapterFound
}
