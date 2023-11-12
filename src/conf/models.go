package conf

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

//  GENERATED CODE — DO NOT EDIT 

var models = []db.Model{
	&m.Adapter{},
	&m.AuditLog{},
	&m.Device{},
	&m.DeviceDiscoveryNotification{},
	&m.DeviceLimitNotification{},
	&m.DeviceStateLog{},
	&m.DeviceStateNotification{},
	&m.Notification{},
	&m.Setting{},
	&m.User{},
	&m.VendorLookup{},
}
