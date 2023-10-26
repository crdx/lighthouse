package users_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/admin/users")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "root")
	assert.Contains(t, res.Body, "anon")
}

func TestViewEditPage(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/admin/users/1/edit")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "root")
}

func TestEditUserWithErrors(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.PostForm("/admin/users/1/edit", map[string]string{
		"username": "",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "required field")
}

func TestChangePassword(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	password := uuid.NewString()

	res := session.PostForm("/admin/users/1/edit", map[string]string{
		"username": "root",
		"password": password,
	})

	assert.Equal(t, 302, res.StatusCode)

	session = helpers.NewSession(auth.StateUnauthenticated)

	res = session.PostForm("/", map[string]string{
		"username": "root",
		"password": password,
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.PostForm("/admin/users/2/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/users")
	assert.NotContains(t, res.Body, "anon")
}

func TestCannotDeleteSelf(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.PostForm("/admin/users/1/delete", nil)
	assert.Equal(t, 400, res.StatusCode)
}
