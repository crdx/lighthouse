package auth_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/adminController"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	helpers.Init()
	app := helpers.App(auth.StateUnauthenticated)
	adminController.InitRoutes(app)
	return helpers.NewSession(app)
}

func TestLoginPage(t *testing.T) {
	session := setup()

	res, body := session.Get("/")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "Username")
	assert.Contains(t, body, "Password")
	assert.NotContains(t, body, "Devices")
}

func TestSuccessfulAdminLogin(t *testing.T) {
	session := setup()

	res, _ := session.PostForm("/", map[string]string{
		"username": "admin",
		"password": "admin",
	})

	assert.Equal(t, 302, res.StatusCode)

	res, body := session.Get("/admin")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "Settings")
}

func TestSuccessfulUserLogin(t *testing.T) {
	session := setup()

	res, body := session.PostForm("/", map[string]string{
		"username": "user",
		"password": "user",
	})

	assert.Equal(t, 302, res.StatusCode)

	res, body = session.Get("/admin")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, body, "Settings")
}

func TestFailedLogin(t *testing.T) {
	session := setup()

	res, body := session.PostForm("/", map[string]string{
		"username": "admin",
		"password": "hunter2",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "Invalid credentials")
	assert.NotContains(t, body, "Devices")
}
