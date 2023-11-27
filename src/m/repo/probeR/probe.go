package probeR

import (
	"time"

	"crdx.org/lighthouse/m/repo/settingR"
)

// TTL returns the maximum duration that a service will be considered to be associated with a device
// once it has been found.
func TTL() time.Duration {
	return 10 * settingR.ServiceScanInterval()
}
