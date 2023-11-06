package vendorLookupR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Unprocessed() []*m.VendorLookup {
	return db.B[m.VendorLookup]("processed = 0").Find()
}
