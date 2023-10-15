package adapterController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/adapterController"
	"crdx.org/lighthouse/controllers/deviceController"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	adapterController.InitRoutes(app)
	deviceController.InitRoutes(app)
	return app
}

func TestEdit(t *testing.T) {
	app := app()

	nameUUID := helpers.UUID()
	vendorUUID := helpers.UUID()

	res, _ := helpers.PostForm(app, "/adapter/1/edit", map[string]string{
		"name":   nameUUID,
		"vendor": vendorUUID,
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := helpers.Get(app, "/device/1")

	assert.Contains(t, body, nameUUID)
	assert.Contains(t, body, vendorUUID)
}

func TestDelete(t *testing.T) {
	app := app()

	_, body := helpers.Get(app, "/device/1/")
	assert.Contains(t, body, "adapter1")

	res, _ := helpers.PostForm(app, "/adapter/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	_, body = helpers.Get(app, "/device/1")
	assert.NotContains(t, body, "adapter1")

	res, _ = helpers.Get(app, "/adapter/1/edit")
	assert.Equal(t, 404, res.StatusCode)
}
