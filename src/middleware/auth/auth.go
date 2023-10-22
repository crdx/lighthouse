package auth

import (
	"net/http"
	"time"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
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

func err(c *fiber.Ctx) error {
	return c.Render("auth/index", fiber.Map{"err": true}, "auth/layout")
}

func needAuth(c *fiber.Ctx) error {
	return c.Render("auth/index", fiber.Map{}, "auth/layout")
}

func logOut(c *fiber.Ctx) error {
	session.Destroy(c)
	return c.Redirect("/")
}

func logIn(c *fiber.Ctx, username string, password string) error {
	user, found := db.B[m.User]().Where("username = ?", username).First()

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

// New returns a fiber.Handler that handles authentication.
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == http.MethodPost && c.Path() == "/goodbye" {
			return logOut(c)
		}

		username := c.FormValue("username")
		password := c.FormValue("password")

		if c.Method() == http.MethodPost && username != "" && password != "" {
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

// AutoLogin returns a fiber.Handler that simulates the user being authorised as the provided state.
// The first user in the db with the required authorisation will be picked. This is designed to be
// used for tests.
func AutoLogin(state State) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := db.B[m.User]().Where("admin = ?", state == StateAdmin).First()
		c.Locals("user", user)
		return c.Next()
	}
}