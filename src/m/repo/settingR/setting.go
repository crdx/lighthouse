package settingR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/setting"
)

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
}

// Invalidate updates the settings cache with current values.
func Invalidate() {
	setting.Update(Map())
}
