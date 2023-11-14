package profileController_test

import (
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/profile")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "admin")
}

func TestChangePassword(t *testing.T) {
	session := helpers.Init(constants.RoleViewer)

	password := uuid.NewString()

	res := session.PostForm("/profile", map[string]string{
		"current_password": "anon",
		"new_password":     password,
	})

	assert.Equal(t, 302, res.StatusCode)

	session = helpers.NewSession(constants.RoleNone)

	res = session.PostForm("/", map[string]string{
		"username": "anon",
		"password": password,
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)
}

func TestCannotChangePasswordWithoutCurrentPassword(t *testing.T) {
	session := helpers.Init(constants.RoleViewer)

	password := uuid.NewString()

	res := session.PostForm("/profile", map[string]string{
		"current_password": "wrong",
		"new_password":     password,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be your current password")
}
