package auth_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestLoginPage(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

	res := session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestSuccessfulAdminLogin(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

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
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleNone)

	res := session.PostForm("/", map[string]string{
		"username": "anon",
		"password": "anon",
		"id":       auth.FormID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "Devices")

	user, _ := db.FindUser(3)
	user.Delete()

	res = session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, auth.FormID)
	assert.NotContains(t, res.Body, "Devices")
}

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		middleware     func(c *fiber.Ctx) error
		role           int64
		expectedStatus int
		roleName       string
	}{
		{auth.Admin, constants.RoleAdmin, 200, "admin"},
		{auth.Admin, constants.RoleEditor, 404, "editor"},
		{auth.Admin, constants.RoleViewer, 404, "viewer"},
		{auth.Editor, constants.RoleAdmin, 200, "admin"},
		{auth.Editor, constants.RoleEditor, 200, "editor"},
		{auth.Editor, constants.RoleViewer, 404, "viewer"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			defer helpers.Start()()
			session := helpers.NewSession(testCase.role, testCase.middleware)
			res := session.Get("/profile")

			assert.Equal(t, testCase.expectedStatus, res.StatusCode)
			if res.StatusCode == 200 { //nolint:usestdlibvars
				assert.Contains(t, res.Body, testCase.roleName)
			} else {
				assert.NotContains(t, res.Body, testCase.roleName)
			}
		})
	}
}
