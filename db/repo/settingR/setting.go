package settingR

import (
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/duration"
	"crdx.org/lighthouse/pkg/util/reflectutil"
)

// Mail

func EnableMail() bool        { return getBool("enable_mail") }
func MailFromHeader() string  { return get("mail_from_header") }
func MailFromAddress() string { return get("mail_from_address") }
func MailToHeader() string    { return get("mail_to_header") }
func MailToAddress() string   { return get("mail_to_address") }
func SMTPHost() string        { return get("smtp_host") }
func SMTPPort() string        { return get("smtp_port") }
func SMTPUser() string        { return get("smtp_user") }
func SMTPPass() string        { return get("smtp_pass") }

// Device Scanning

func EnableDeviceScan() bool            { return getBool("enable_device_scan") }
func DeviceScanInterval() time.Duration { return getDuration("device_scan_interval") }
func NotifyOnNewDevice() bool           { return getBool("notify_on_new_device") }
func WatchNew() bool                    { return getBool("watch_new") }
func PingNew() bool                     { return getBool("ping_new") }

// Service Scanning

func EnableServiceScan() bool            { return getBool("enable_service_scan") }
func ServiceScanInterval() time.Duration { return getDuration("service_scan_interval") }
func NotifyOnNewService() bool           { return getBool("notify_on_new_service") }

// System

func MACVendorsAPIKey() string { return get("macvendors_api_key") }
func Timezone() string         { return get("timezone") }

// Other

func SourceMACAddresses() string { return get("source_mac_addresses") }

var cache map[string]string

// invalidate invalidates the settings cache.
func invalidate() {
	cache = nil
}

// Map returns all settings as a map[string]string.
func Map() map[string]string {
	return db.MapBy2[string, string]("Name", "Value", db.FindSettings())
}

// Set sets a setting to value.
func Set(name string, value any) {
	s := reflectutil.ToString(value)

	setting, found := db.FindSettingByName(name)
	if !found {
		db.CreateSetting(&db.Setting{
			Name:  name,
			Value: s,
		})
	} else {
		setting.UpdateValue(s)
	}
	invalidate()
}

// get returns a setting as a string.
func get(name string) string {
	if cache == nil {
		cache = Map()
	}

	return cache[name]
}

// getDuration returns a setting as a time.Duration.
func getDuration(name string) time.Duration {
	return duration.MustParse(get(name))
}

// // getInt returns a setting as an int.
// func getInt(name string) int {
// 	return int(lo.Must(strconv.ParseInt(get(name), 10, 64)))
// }

// // getFloat returns a setting as a float.
// func getFloat(name string) float64 {
// 	return lo.Must(strconv.ParseFloat(get(name), 64))
// }

// getBool returns a setting as a bool.
func getBool(name string) bool {
	return get(name) == "1"
}
