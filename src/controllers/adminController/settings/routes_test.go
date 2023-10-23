package settings_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateAdmin)
}

func TestList(t *testing.T) {
	session := setup()
	res := session.Get("/admin/settings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "MACVendors")
}

func TestEdit(t *testing.T) {
	session := setup()

	apiKey := uuid.NewString()

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")

	assert.Contains(t, res.Body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	session := setup()

	id := uuid.NewString()
	apiKey := strings.Repeat(id, 100)

	res := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, id)
	assert.Contains(t, res.Body, "must be a maximum of")
	assert.Contains(t, res.Body, "characters in length")
}
