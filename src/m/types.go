package m

import (
	"crdx.org/lighthouse/models/adapterM"
	"crdx.org/lighthouse/models/deviceM"
	"crdx.org/lighthouse/models/deviceStateLogM"
	"crdx.org/lighthouse/models/networkM"
	"crdx.org/lighthouse/models/notificationM"
)

type (
	Network        = networkM.Network
	Device         = deviceM.Device
	Adapter        = adapterM.Adapter
	DeviceStateLog = deviceStateLogM.DeviceStateLog
	Notification   = notificationM.Notification
)

var All = []any{
	&Network{},
	&Device{},
	&Adapter{},
	&DeviceStateLog{},
	&Notification{},
}
