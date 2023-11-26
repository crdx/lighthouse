package probeR

import (
	"time"

	"crdx.org/lighthouse/m/repo/settingR"
)

func TTL() time.Duration {
	return 10 * settingR.ServiceScanInterval()
}
