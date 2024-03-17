package db

// IsOnline returns true if this adapter was last seen within its grace period.
func IsOnline(device *Device, adapter *Adapter) bool {
	return adapter.LastSeenAt.After(Now().Add(-device.GracePeriodDuration()))
}

// IsNotResponding returns true if this adapter was last seen within half of the grace
// period. This indicates a device that may be about to go offline.
func IsNotResponding(device *Device, adapter *Adapter) bool {
	return adapter.LastSeenAt.Before(Now().Add(-device.GracePeriodDuration() / 2))
}
