package adminController_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/controllers/adminController"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	helpers.Init()
	app := helpers.App(auth.StateAdmin)
	adminController.InitRoutes(app)
	return helpers.NewSession(app)
}

func TestIndex(t *testing.T) {
	session := setup()
	res, body := session.Get("/admin")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "MACVendors")
}

func TestSave(t *testing.T) {
	session := setup()

	apiKey := helpers.UUID()

	res, _ := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
		"timezone":           "Europe/London",
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := session.Get("/admin")

	assert.Contains(t, body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	session := setup()

	apiKey := strings.Repeat(helpers.UUID(), 100)

	res, body := session.PostForm("/admin/settings", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "must be a maximum of")
	assert.Contains(t, body, "characters in length")
}
