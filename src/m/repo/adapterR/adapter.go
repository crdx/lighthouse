package adapterR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func FindByMACAddress(macAddress string) (*m.Adapter, bool) {
	return db.B[m.Adapter]("mac_address = ?", macAddress).First()
}
