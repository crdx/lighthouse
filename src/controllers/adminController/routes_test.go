package adminController_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/controllers/adminController"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	adminController.InitRoutes(app)
	return app
}

func TestIndex(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/admin")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "MACVendors")
}

func TestSave(t *testing.T) {
	app := app()

	apiKey := helpers.UUID()

	res, _ := helpers.PostForm(app, "/admin", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := helpers.Get(app, "/admin")

	assert.Contains(t, body, apiKey)
}

func TestEditWithErrors(t *testing.T) {
	app := app()

	apiKey := strings.Repeat(helpers.UUID(), 100)

	res, body := helpers.PostForm(app, "/admin", map[string]string{
		"macvendors_api_key": apiKey,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "must be a maximum of")
	assert.Contains(t, body, "characters in length")
}
