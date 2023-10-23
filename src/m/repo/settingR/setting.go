package settingR

import (
	"strconv"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"github.com/samber/lo"
)

const (
	Watch = "watch"

	EnableMail      = "enable_mail"
	MailFromHeader  = "mail_from_header"
	MailFromAddress = "mail_from_address"
	MailToHeader    = "mail_to_header"
	MailToAddress   = "mail_to_address"
	SMTPHost        = "smtp_host"
	SMTPPort        = "smtp_port"
	SMTPUser        = "smtp_user"
	SMTPPass        = "smtp_pass"

	MACVendorsAPIKey = "macvendors_api_key"
	Timezone         = "timezone"
)

var cache map[string]string

// Invalidate invalidates the settings cache.
func Invalidate() {
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
	Invalidate()
}

// Get returns a setting as a string.
func Get(name string) string {
	if cache == nil {
		cache = Map()
	}

	return cache[name]
}

// Get returns a setting as an int.
func GetInt(name string) int {
	return int(lo.Must(strconv.ParseInt(Get(name), 10, 64)))
}

// Get returns a setting as a uint.
func GetUint(name string) uint {
	return uint(lo.Must(strconv.ParseUint(Get(name), 10, 64)))
}

// Get returns a setting as a float.
func GetFloat(name string) float64 {
	return lo.Must(strconv.ParseFloat(Get(name), 64))
}

// Get returns a setting as a bool.
func GetBool(name string) bool {
	return Get(name) == "1"
}
