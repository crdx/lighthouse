package deviceController_test

import (
	"testing"

	"crdx.org/db"
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

func TestList(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "device1")
}

func TestView(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/device/1")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, body, "127.0.0.1")
	assert.Contains(t, body, "adapter1")
	assert.Contains(t, body, "Corp 1")
}

func TestEdit(t *testing.T) {
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

func TestEditWithErrors(t *testing.T) {
	app := app()

	nameUUID := helpers.UUID()
	notesUUID := helpers.UUID()
	iconUUID := helpers.UUID()

	res, body := helpers.PostForm(app, "/device/1/edit", map[string]string{
		"name":         nameUUID,
		"notes":        notesUUID,
		"icon":         iconUUID,
		"grace_period": "",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "required field")

	_, body = helpers.Get(app, "/device/1")

	assert.NotContains(t, body, notesUUID)
	assert.NotContains(t, body, iconUUID)
}

func TestMerge(t *testing.T) {
	app := app()

	res, _ := helpers.PostForm(app, "/device/1/merge", map[string]string{
		"device_id": "2",
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := helpers.Get(app, "/device/1")

	assert.Contains(t, body, "2023-10-01")
	assert.Contains(t, body, "adapter1")
	assert.Contains(t, body, "adapter2")

	device := lo.Must(db.First[m.Device](1))

	assert.Len(t, device.Adapters(), 2)
	assert.NotNil(t, device.DeletedAt)

	_, found := db.First[m.Device](2)
	assert.False(t, found)
}

func TestDelete(t *testing.T) {
	app := app()

	res, _ := helpers.Get(app, "/device/1")
	assert.Equal(t, 200, res.StatusCode)

	res, _ = helpers.PostForm(app, "/device/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res, _ = helpers.Get(app, "/device/1")
	assert.Equal(t, 404, res.StatusCode)
}
