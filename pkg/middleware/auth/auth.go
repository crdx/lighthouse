package auth

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/session"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

// FormID is used to make sure with no shadow of a doubt that the submitted form is the login form.
const FormID = "afc434ce-bf57-48f7-9844-e9ab4091f19a"

func err(c fiber.Ctx) error {
	time.Sleep(100 * time.Millisecond)
	return c.Render("auth/index", fiber.Map{"err": true, "id": FormID}, "auth/layout")
}

func needAuth(c fiber.Ctx) error {
	return c.Render("auth/index", fiber.Map{"id": FormID}, "auth/layout")
}

func logOut(c fiber.Ctx) error {
	auditLogR.Add(c, "User %s logged out", globals.CurrentUser(c).Username)
	session.Destroy(c)
	return c.Redirect().To("/")
}

func logIn(c fiber.Ctx, username string, password string) error {
	user, found := db.FindUserByUsername(username)

	if !found {
		auditLogR.Add(c, "Unknown user %q tried to log in", username)
		return err(c)
	}

	if !stringutil.VerifyHashAndPassword(user.PasswordHash, password) {
		auditLogR.Add(c, "User %s failed to log in", username)
		return err(c)
	}

	user.UpdateLastLoginAt(db.Now())
	auditLogR.Add(c, "User %s logged in", user.Username)

	session.Set(c, "user_id", user.ID)
	return c.Redirect().To(c.Path())
}

// New returns middleware that handles authentication.
func New() fiber.Handler {
	return func(c fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		isAuthForm := c.FormValue("id") == FormID

		if isAuthForm && c.Method() == http.MethodPost && username != "" && password != "" {
			return logIn(c, username, password)
		}

		userId := session.Get[int64](c, "user_id")
		if userId == 0 {
			return needAuth(c)
		}

		user, found := db.FindUser(userId)
		if !found {
			return needAuth(c)
		}

		user.UpdateLastVisitAt(db.Now())
		c.Locals(globals.CurrentUserKey, user)

		if c.Method() == http.MethodPost && c.Path() == "/bye" {
			return logOut(c)
		}

		return c.Next()
	}
}

// Admin is middleware that only allows the request to continue if the current user is an admin.
func Admin(c fiber.Ctx) error {
	if globals.CurrentUser(c).IsAdmin() {
		return c.Next()
	}
	return c.SendStatus(404)
}

// Editor is middleware that only allows the request to continue if the current user is an editor.
func Editor(c fiber.Ctx) error {
	if globals.CurrentUser(c).IsEditor() {
		return c.Next()
	}

	return c.SendStatus(404)
}

// AutoLogin returns middleware that simulates the user being authorised as the provided state. The
// first user in the db with the required authorisation will be picked.
func AutoLogin(role int64) fiber.Handler {
	return func(c fiber.Ctx) error {
		if session.Get[int64](c, "user_id") == 0 {
			user, _ := db.FindUserByRole(role)
			user.UpdateLastLoginAt(db.Now())
			auditLogR.Add(c, "User %s logged in", user.Username)
			session.Set(c, "user_id", user.ID)
		}

		userId := session.Get[int64](c, "user_id")
		user, _ := db.FindUser(userId)
		user.UpdateLastVisitAt(db.Now())
		c.Locals(globals.CurrentUserKey, user)
		return c.Next()
	}
}

func ExpireUserID(userId int64, exceptSessionId string) {
	sessions := db.FindOtherSessions(exceptSessionId)

	for _, session := range sessions {
		sessionData := map[any]any{}
		lo.Must0(gob.NewDecoder(bytes.NewBuffer(session.V)).Decode(&sessionData))

		if sessionUserId := sessionData["user_id"].(int64); sessionUserId == userId {
			db.DeleteSession(session.K)
		}
	}
}
