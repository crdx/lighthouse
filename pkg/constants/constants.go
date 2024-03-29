package constants

const UnknownVendorLabel = "Unknown"

const (
	DefaultDeviceIconClass = "duotone:question"
	DefaultGracePeriod     = "5 mins"
)

const (
	WatchColumnLabel = `<span class="icon is-small"><i class="fa-duotone fa-eye"></i></span>`
	TypeColumnLabel  = `<span class="icon"><i class="fa-duotone fa-exchange-alt fa-sm"></i></span>`
)

const (
	ActivityRowsPerPage     = 100
	NotificationRowsPerPage = 100
)

const (
	TimeFormatReadablePrecise = "15:04:05 on Mon, Jan 2 2006"
	TimeFormatReadable        = "15:04 on Mon, Jan 2 2006"
	TimeFormatSystem          = "2006-01-02 15:04:05 MST"
	TimeFormatEuropeanKitchen = "15:04"
)

const (
	RoleNone   int64 = 0
	RoleViewer int64 = 1
	RoleEditor int64 = 2
	RoleAdmin  int64 = 3
)
