package setting

import (
	"strconv"

	"github.com/samber/lo"
)

var cache map[string]string

const Watch = "watch"
const MACVendorsAPIKey = "macvendors_api_key"
const Timezone = "timezone"

// Update updates the settings cache.
func Update(settings map[string]string) {
	cache = settings
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

// Get returns a setting as a string.
func Get(name string) string {
	return cache[name]
}
