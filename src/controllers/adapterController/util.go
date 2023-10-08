package adapterController

import (
	"strconv"

	"crdx.org/lighthouse/m"
)

func getAdapter(idStr string) (*m.Adapter, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, false
	}

	return m.ForAdapter(uint(id)).First()
}
