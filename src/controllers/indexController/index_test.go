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
