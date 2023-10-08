package indexController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/indexController"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	indexController.InitRoutes(app)
	return app
}

func TestDeviceList(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:BB:CC:DD:EE:FF")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "localhost")
}

func TestViewDevice(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/device/1")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:BB:CC:DD:EE:FF")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "localhost")
	assert.Contains(t, body, "adapter1")
	assert.Contains(t, body, "Computer Corporation")
}

func TestEditDevice(t *testing.T) {
	app := app()

	nameUUID := helpers.UUID()
	notesUUID := helpers.UUID()
	iconUUID := helpers.UUID()

	res, _ := helpers.PostForm(app, "/device/1/edit", map[string]string{
		"name":         nameUUID,
		"notes":        notesUUID,
		"icon":         iconUUID,
		"grace_period": "6",
	})

	assert.Equal(t, 302, res.StatusCode)

	res, body := helpers.Get(app, "/device/1")

	assert.Contains(t, body, nameUUID)
	assert.Contains(t, body, notesUUID)
	assert.Contains(t, body, iconUUID)
}
