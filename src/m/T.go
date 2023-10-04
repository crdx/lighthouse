package m

import "crdx.org/db"

var All = []any{
	&Network{},
	&Device{},
	&Adapter{},
	&DeviceStateLog{},
}

// ForAdapter returns a db.Builder for the Adapter with the specified ID.
func ForAdapter(id uint) *db.Builder[Adapter] {
	return db.B(Adapter{ID: id})
}

// ForDevice returns a db.Builder for the Device with the specified ID.
func ForDevice(id uint) *db.Builder[Device] {
	return db.B(Device{ID: id})
}
