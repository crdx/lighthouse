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
	// Cannot query via m.VendorLookup{Processed: false} because false is the zero value of bool so
	// it's ignored.
	return db.B[m.VendorLookup]().Where("processed = 0").Find()
}
