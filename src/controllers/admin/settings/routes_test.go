package settings_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestList(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/admin/settings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "MACVendors")
	assert.Contains(t, res.Body, "Settings")
}

func TestViewerCannotList(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleViewer)

	res := session.Get("/admin/settings")
	assert.Equal(t, 404, res.StatusCode)
}

func TestEdit(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	apiKey := uuid.NewString()

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key":    apiKey,
		"timezone":              "Europe/London",
		"device_scan_interval":  "1 min",
		"service_scan_interval": "1 hour",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")
	assert.Contains(t, res.Body, apiKey)
}

func TestViewerCannotEdit(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleViewer)

	apiKey := uuid.NewString()

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
	})

	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	apiKey := strings.Repeat(uuid.NewString(), 20)

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, apiKey)
	assert.Contains(t, res.Body, "must be a maximum of")
	assert.Contains(t, res.Body, "characters in length")
}

func TestCacheInvalidation(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	currentTimezone := settingR.Timezone()

	session.PostForm("/admin/settings", map[string]string{
		"timezone":              "America/New_York",
		"device_scan_interval":  "1 min",
		"service_scan_interval": "1 hour",
	})

	timezone := settingR.Timezone()

	assert.NotEqual(t, timezone, currentTimezone)
	assert.Equal(t, "America/New_York", timezone)
}
