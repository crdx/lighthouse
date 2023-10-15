package activityController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/activityController"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	activityController.InitRoutes(app)
	return app
}

func TestList(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/activity")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "device1")
	assert.Contains(t, body, "device2")
	assert.NotContains(t, body, "device3")
	assert.Contains(t, body, "online")
	assert.Contains(t, body, "offline")
}
