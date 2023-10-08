package deviceController

import (
	"strconv"

	"crdx.org/lighthouse/m"
)

func getDevice(idStr string) (*m.Device, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, false
	}

	return m.ForDevice(uint(id)).First()
}
