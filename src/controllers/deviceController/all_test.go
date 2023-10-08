package deviceController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/deviceController"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	deviceController.InitRoutes(app)
	return app
}

func TestDeviceList(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "device1")
}

func TestViewDevice(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/device/1")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "adapter1")
	assert.Contains(t, body, "Corp 1")
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

	_, body := helpers.Get(app, "/device/1")

	assert.Contains(t, body, nameUUID)
	assert.Contains(t, body, notesUUID)
	assert.Contains(t, body, iconUUID)
}

func TestEditDeviceWithErrors(t *testing.T) {
	app := app()

	notesUUID := helpers.UUID()
	iconUUID := helpers.UUID()

	res, body := helpers.PostForm(app, "/device/1/edit", map[string]string{
		"name":         "",
		"notes":        notesUUID,
		"icon":         iconUUID,
		"grace_period": "6",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "Name is a required field")

	_, body = helpers.Get(app, "/device/1")

	assert.NotContains(t, body, notesUUID)
	assert.NotContains(t, body, iconUUID)
}

func TestMergeDevice(t *testing.T) {
	app := app()

	res, _ := helpers.PostForm(app, "/device/1/merge", map[string]string{
		"device_id": "2",
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := helpers.Get(app, "/device/1")

	assert.Contains(t, body, "01/10/2023")
	assert.Contains(t, body, "adapter1")
	assert.Contains(t, body, "adapter2")

	device := lo.Must(m.ForDevice(1).First())

	assert.Len(t, device.Adapters(), 2)
	assert.NotNil(t, device.DeletedAt)

	_, found := m.ForDevice(2).First()
	assert.False(t, found)
}
