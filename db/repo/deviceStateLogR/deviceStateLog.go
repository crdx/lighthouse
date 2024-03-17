package deviceStateLogR

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/pager"
)

func GetListRowCount(deviceID int64) int {
	if deviceID != 0 {
		return int(db.CountDeviceStateLogsListViewForDevice(deviceID))
	} else {
		return int(db.CountDeviceStateLogsListView())
	}
}

func GetList(page int, perPage int, deviceID int64) []*db.DeviceStateLogsView {
	offset := int64(pager.GetOffset(page, perPage))
	limit := int64(perPage)

	if deviceID != 0 {
		return db.FindDeviceStateLogsListViewForDevice(deviceID, offset, limit)
	} else {
		return db.FindDeviceStateLogsListView(offset, limit)
	}
}
