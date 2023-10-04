package deviceStateLogR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func All() []*m.DeviceStateLog {
	return db.B[m.DeviceStateLog]().Find()
}

// —————————————————————————————————————————————————————————————————————————————————————————————————
