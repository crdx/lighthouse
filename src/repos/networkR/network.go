package networkR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func All() []*m.Network {
	return db.B[m.Network]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func Upsert(ipRange string) (*m.Network, bool) {
	network, found := db.FirstOrInit(m.Network{IPRange: ipRange})

	if !found {
		network.Name = ipRange
		network.AlertOnNewDevices = true

		db.Save(&network)
	}

	return network, found
}
