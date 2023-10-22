package vendorLookupR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func All() []*m.VendorLookup {
	return db.B[m.VendorLookup]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func Unprocessed() []*m.VendorLookup {
	return db.B[m.VendorLookup]("processed = 0").Find()
}
