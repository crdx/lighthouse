package notificationController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/notificationController"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func app() *fiber.App {
	helpers.Init()
	app := helpers.App()
	notificationController.InitRoutes(app)
	return app
}

func TestList(t *testing.T) {
	app := app()
	res, body := helpers.Get(app, "/notifications")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "a thing has happened")
	assert.Contains(t, body, "here are more details about the thing that happened")
}
