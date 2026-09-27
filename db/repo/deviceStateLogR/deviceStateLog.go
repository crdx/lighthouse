package deviceStateLogR

import (
	"fmt"

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

func GetList(page int, perPage int, deviceID int64, sortColumn string, sortDirection string) []*db.DeviceStateLogsView {
	orderByTemplates := map[string]string{
		"type":  "icon %s, created_at DESC, device_id ASC",
		"name":  "name %s, created_at DESC, device_id ASC",
		"state": "state %s, created_at DESC, device_id ASC",
		"when":  "created_at %s, device_id ASC",
		"date":  "created_at %s, device_id ASC",
	}
	sortDirections := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	orderByTemplate, columnOK := orderByTemplates[sortColumn]
	direction, directionOK := sortDirections[sortDirection]
	if !columnOK || !directionOK {
		return nil
	}

	offset := int64(pager.GetOffset(page, perPage))
	limit := int64(perPage)
	orderBy := fmt.Sprintf(orderByTemplate, direction)

	if deviceID != 0 {
		return db.ScanN[db.DeviceStateLogsView](fmt.Sprintf(`
			SELECT created_at, device_id, name, icon, deleted_at, state
			FROM device_state_logs_view
			WHERE device_id = ?
			ORDER BY %s
			LIMIT ?, ?
		`, orderBy), deviceID, offset, limit)
	}

	return db.ScanN[db.DeviceStateLogsView](fmt.Sprintf(`
		SELECT created_at, device_id, name, icon, deleted_at, state
		FROM device_state_logs_view
		ORDER BY %s
		LIMIT ?, ?
	`, orderBy), offset, limit)
}
