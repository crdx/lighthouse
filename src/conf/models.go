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
	&m.DeviceIPAddressLog{},
	&m.DeviceLimitNotification{},
	&m.DeviceServiceNotification{},
	&m.DeviceStateLog{},
	&m.DeviceStateNotification{},
	&m.Mapping{},
	&m.Notification{},
	&m.Scan{},
	&m.ScanResult{},
	&m.Service{},
	&m.Setting{},
	&m.User{},
	&m.VendorLookup{},
}
