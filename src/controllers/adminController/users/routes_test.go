package users_test

import (
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/admin/users")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "root")
	assert.Contains(t, res.Body, "anon")
}

func TestViewerCannotList(t *testing.T) {
	session := helpers.Init(constants.RoleViewer)

	res := session.Get("/admin/users")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewEditPage(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/admin/users/1/edit")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "root")
}

func TestCannotEditUsername(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/1/edit", map[string]string{
		"username": "joe",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/users")
	assert.Contains(t, res.Body, "root")
}

func TestEditWithErrors(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/1/edit", map[string]string{
		"password": "foo",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be at least 4 characters in length")
}

func TestCreate(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	password := uuid.NewString()

	res := session.PostForm("/admin/users/create", map[string]string{
		"username": "joe",
		"password": password,
		"role":     "1",
	})

	assert.Equal(t, 302, res.StatusCode)

	session = helpers.NewSession(constants.RoleNone)

	res = session.PostForm("/", map[string]string{
		"username": "joe",
		"password": password,
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)
}

func TestViewCreatePage(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/admin/users/create")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Create User")
}

func TestCannotCreateWithUnavailableUsername(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/create", map[string]string{
		"username": "root",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be an available username")
}

func TestCannotCreateWithInvalidRole(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/create", map[string]string{
		"role": "100",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be a valid role")
}

func TestChangePassword(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	password := uuid.NewString()

	res := session.PostForm("/admin/users/1/edit", map[string]string{
		"username": "root",
		"password": password,
	})

	assert.Equal(t, 302, res.StatusCode)

	session = helpers.NewSession(constants.RoleNone)

	res = session.PostForm("/", map[string]string{
		"username": "root",
		"password": password,
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/3/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/users")
	assert.NotContains(t, res.Body, "anon")
}

func TestCannotDeleteSelf(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.PostForm("/admin/users/1/delete", nil)
	assert.Equal(t, 400, res.StatusCode)
}
