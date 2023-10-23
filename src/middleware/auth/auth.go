package auth

import (
	"net/http"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/util/stringutil"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

type State int

const (
	StateAdmin State = iota
	StateUser
	StateUnauthenticated
)

// Used to make sure with no shadow of a doubt that the submitted form is the login form.
const ID = "afc434ce-bf57-48f7-9844-e9ab4091f19a"

func err(c *fiber.Ctx) error {
	time.Sleep(100 * time.Millisecond)
	return c.Render("auth/index", fiber.Map{"err": true, "id": ID}, "auth/layout")
}

func needAuth(c *fiber.Ctx) error {
	return c.Render("auth/index", fiber.Map{"id": ID}, "auth/layout")
}

func logOut(c *fiber.Ctx) error {
	session.Destroy(c)
	return c.Redirect("/")
}

func logIn(c *fiber.Ctx, username string, password string) error {
	user, found := db.B[m.User]("username = ?", username).First()

	if !found {
		return err(c)
	}

	if !stringutil.VerifyHashAndPassword(user.PasswordHash, password) {
		return err(c)
	}

	user.Update("last_login", time.Now())

	session.Set(c, "user_id", user.ID)
	return c.Redirect(c.Path())
}

// New returns middleware that handles authentication.
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == http.MethodPost && c.Path() == "/goodbye" {
			return logOut(c)
		}

		username := c.FormValue("username")
		password := c.FormValue("password")
		isAuthForm := c.FormValue("id") == ID

		if isAuthForm && c.Method() == http.MethodPost && username != "" && password != "" {
			return logIn(c, username, password)
		}

		userId := session.GetUint(c, "user_id")
		if userId == 0 {
			return needAuth(c)
		}

		user, found := db.First[m.User](userId)
		if !found {
			return needAuth(c)
		}

		c.Locals("user", user)
		return c.Next()
	}
}

// Admin is middleware that only allows the request to continue if the current user is an admin.
func Admin(c *fiber.Ctx) error {
	if !globals.IsAdmin(c) {
		return c.SendStatus(404)
	}
	return c.Next()
}

// AutoLogin returns middleware that simulates the user being authorised as the provided state. The
// first user in the db with the required authorisation will be picked. This is designed to be used
// for tests.
func AutoLogin(state State) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := db.B[m.User]("admin = ?", state == StateAdmin).First()
		c.Locals("user", user)
		return c.Next()
	}
}
