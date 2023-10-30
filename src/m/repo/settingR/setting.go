package settingR

import (
	"strconv"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"github.com/samber/lo"
)

func Watch() bool              { return getBool("watch") }
func EnableMail() bool         { return getBool("enable_mail") }
func MailFromHeader() string   { return get("mail_from_header") }
func MailFromAddress() string  { return get("mail_from_address") }
func MailToHeader() string     { return get("mail_to_header") }
func MailToAddress() string    { return get("mail_to_address") }
func SMTPHost() string         { return get("smtp_host") }
func SMTPPort() string         { return get("smtp_port") }
func SMTPUser() string         { return get("smtp_user") }
func SMTPPass() string         { return get("smtp_pass") }
func MACVendorsAPIKey() string { return get("macvendors_api_key") }
func Timezone() string         { return get("timezone") }
func Passive() bool            { return getBool("passive") }

var cache map[string]string

// invalidate invalidates the settings cache.
func invalidate() {
	cache = nil
}

// Map returns all settings as a map[string]string.
func Map() map[string]string {
	settings := map[string]string{}

	for _, setting := range db.B[m.Setting]().Find() {
		settings[setting.Name] = setting.Value
	}

	return settings
}

// Set sets a setting to value.
func Set(name string, value any) {
	setting, _ := db.FirstOrCreate(m.Setting{Name: name})
	setting.Update("value", value)
	invalidate()
}

// Get returns a setting as a string.
func get(name string) string {
	if cache == nil {
		cache = Map()
	}

	return cache[name]
}

// // Get returns a setting as an int.
// func getInt(name string) int {
// 	return int(lo.Must(strconv.ParseInt(get(name), 10, 64)))
// }

// // Get returns a setting as a uint.
// func getUint(name string) uint {
// 	return uint(lo.Must(strconv.ParseUint(get(name), 10, 64)))
// }

// // Get returns a setting as a float.
// func getFloat(name string) float64 {
// 	return lo.Must(strconv.ParseFloat(get(name), 64))
// }

// Get returns a setting as a bool.
func getBool(name string) bool {
	return get(name) == "1"
}
