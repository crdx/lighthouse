package globals

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

type Values struct {
	Flash        *flash.Message
	User         *m.User
	CurrentRoute string
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
	values := Values{
		CurrentRoute: c.Route().Name,
	}

	if flashMessage, found := session.GetOnce[flash.Message](c, "globals.flash"); found {
		values.Flash = &flashMessage
	}

	values.User = CurrentUser(c)

	return &values
}
