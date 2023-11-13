package globals

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/userR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

type Values struct {
	Flash *flash.Message
	User  *m.User

	UserIsViewer bool
	UserIsEditor bool
	UserIsAdmin  bool
}

const CurrentUserKey = "globals.current_user"

// CurrentUser returns the current user from the session.
func CurrentUser(c *fiber.Ctx) *m.User {
	user, _ := c.Locals(CurrentUserKey).(*m.User)
	return user
}

// IsCurrentUser returns whether user is the current user.
func IsCurrentUser(c *fiber.Ctx, user *m.User) bool {
	return CurrentUser(c).ID == user.ID
}

// Get returns the encapsulated globals to be referenced from templates.
func Get(c *fiber.Ctx) *Values {
	values := Values{}

	if flashMessage, found := session.GetOnce[flash.Message](c, "globals.flash"); found {
		values.Flash = &flashMessage
	}

	user := CurrentUser(c)
	values.User = user

	if user != nil {
		values.UserIsAdmin = user.Role >= userR.RoleAdmin
		values.UserIsEditor = user.Role >= userR.RoleEditor
		values.UserIsViewer = user.Role >= userR.RoleViewer
	}

	return &values
}
