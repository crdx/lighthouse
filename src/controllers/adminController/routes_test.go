package adminController_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateAdmin)
}

func TestListSettings(t *testing.T) {
	session := setup()
	res, body := session.Get("/admin/settings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "MACVendors")
}

func TestSaveSettings(t *testing.T) {
	session := setup()

	apiKey := stringutil.UUID()

	res, _ := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := session.Get("/admin/settings")

	assert.Contains(t, body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	session := setup()

	apiKey := strings.Repeat(stringutil.UUID(), 100)

	res, body := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "must be a maximum of")
	assert.Contains(t, body, "characters in length")
}
