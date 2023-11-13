package auth_test

import (
	"testing"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestLoginPage(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestSuccessfulAdminLogin(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "root",
		"password": "root",
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Settings")
	assert.NotContains(t, res.Body, auth.FormID)
}

func TestSuccessfulUserLogin(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "anon",
		"password": "anon",
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/settings")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, "Settings")
}

func TestInvalidUsername(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "john",
		"password": "root",
		"id":       auth.FormID,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Invalid credentials")
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestInvalidPassword(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "root",
		"password": "hunter2",
		"id":       auth.FormID,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Invalid credentials")
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestInvalidFormID(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "root",
		"password": "root",
		"id":       uuid.NewString(),
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestLogout(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	session.PostForm("/", map[string]string{
		"username": "root",
		"password": "root",
		"id":       auth.FormID,
	})

	res := session.PostForm("/bye", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/")
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestUserIsDeletedWhileLoggedIn(t *testing.T) {
	session := helpers.Init(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "anon",
		"password": "anon",
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Devices")

	lo.Must(db.First[m.User](3)).Delete()

	res = session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}
