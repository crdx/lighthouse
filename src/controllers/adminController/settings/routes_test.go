package settings_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/admin/settings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "MACVendors")
	assert.Contains(t, res.Body, "Settings")
}

func TestUserCannotList(t *testing.T) {
	session := helpers.Init(auth.StateUser)

	res := session.Get("/admin/settings")
	assert.Equal(t, 404, res.StatusCode)
}

func TestEdit(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	apiKey := uuid.NewString()

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
		"scan_interval":      "1 min",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")
	assert.Contains(t, res.Body, apiKey)
}

func TestUserCannotEdit(t *testing.T) {
	session := helpers.Init(auth.StateUser)

	apiKey := uuid.NewString()

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
	})

	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

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
	session := helpers.Init(auth.StateAdmin)

	currentTimezone := settingR.Timezone()

	session.PostForm("/admin/settings", map[string]string{
		"timezone":      "America/New_York",
		"scan_interval": "1 min",
	})

	timezone := settingR.Timezone()

	assert.NotEqual(t, timezone, currentTimezone)
	assert.Equal(t, "America/New_York", timezone)
}
