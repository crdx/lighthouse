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

func TestEditAdapter(t *testing.T) {
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
