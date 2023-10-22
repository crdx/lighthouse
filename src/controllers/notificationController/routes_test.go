package notificationController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/notificationController"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	helpers.Init()
	app := helpers.App(auth.StateAdmin)
	notificationController.InitRoutes(app)
	return helpers.NewSession(app)
}

func TestList(t *testing.T) {
	session := setup()
	res, body := session.Get("/notifications")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "a thing has happened")
	assert.Contains(t, body, "here are more details about the thing that happened")
}
