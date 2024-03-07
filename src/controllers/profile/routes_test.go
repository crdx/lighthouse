package profile_test

import (
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestView(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/profile")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "admin")
}

func TestChangePassword(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleEditor)

	password := uuid.NewString()

	res := session.PostForm("/profile", map[string]string{
		"current_password":     "ed",
		"new_password":         password,
		"confirm_new_password": password,
	})

	assert.Equal(t, 302, res.StatusCode)

	session = helpers.NewSession(constants.RoleNone)

	res = session.PostForm("/", map[string]string{
		"username": "ed",
		"password": password,
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)
}

func TestCannotChangePasswordWithoutCurrentPassword(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleEditor)

	password := uuid.NewString()

	res := session.PostForm("/profile", map[string]string{
		"current_password":     "wrong",
		"new_password":         password,
		"confirm_new_password": password,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be your current password")
}

func TestCannotChangePasswordWithoutMatchingPasswords(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/profile", map[string]string{
		"current_password":     "root",
		"new_password":         "hunter2",
		"confirm_new_password": "hunter3",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "passwords must match")
}
