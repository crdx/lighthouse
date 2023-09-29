package m

import (
	"crdx.org/lighthouse/models/deviceMappingModel"
	"crdx.org/lighthouse/models/deviceModel"
	"crdx.org/lighthouse/models/deviceStateLogModel"
	"crdx.org/lighthouse/models/networkModel"
	"crdx.org/lighthouse/models/notificationModel"
)

type (
	Network        = networkModel.Network
	Device         = deviceModel.Device
	DeviceMapping  = deviceMappingModel.DeviceMapping
	DeviceStateLog = deviceStateLogModel.DeviceStateLog
	Notification   = notificationModel.Notification
)

var All = []any{
	&Network{},
	&Device{},
	&DeviceMapping{},
	&DeviceStateLog{},
	&Notification{},
}
