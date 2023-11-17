package mappingR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Map() map[string]string {
	mappings := map[string]string{}
	for _, mapping := range db.B[m.Mapping]().Find() {
		mappings[mapping.IPAddress] = mapping.MACAddress
	}
	return mappings
}
