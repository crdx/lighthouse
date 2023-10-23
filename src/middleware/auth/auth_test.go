package auth_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateUnauthenticated)
}

func TestLoginPage(t *testing.T) {
	session := setup()

	res := session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Username")
	assert.Contains(t, res.Body, "Password")
	assert.NotContains(t, res.Body, "Devices")
}

func TestSuccessfulAdminLogin(t *testing.T) {
	session := setup()

	res := session.PostForm("/", map[string]string{
		"username": "admin",
		"password": "admin",
		"id":       auth.ID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Settings")
}

func TestSuccessfulUserLogin(t *testing.T) {
	session := setup()

	res := session.PostForm("/", map[string]string{
		"username": "user",
		"password": "user",
		"id":       auth.ID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, "Settings")
}

func TestFailedLogin(t *testing.T) {
	session := setup()

	res := session.PostForm("/", map[string]string{
		"username": "admin",
		"password": "hunter2",
		"id":       auth.ID,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Invalid credentials")
	assert.NotContains(t, res.Body, "Devices")
}
