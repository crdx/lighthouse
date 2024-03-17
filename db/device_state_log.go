package db

import (
	"crdx.org/lighthouse/pkg/util"
)

func (self *DeviceStateLogsView) IconClass() string {
	return util.IconToClass(self.Icon)
}
