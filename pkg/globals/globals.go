package globals

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session/v3"
	"github.com/gofiber/fiber/v3"
)

type Values struct {
	Flash        *flash.Message
	User         *db.User
	CurrentRoute string
}

const CurrentUserKey = "globals.current_user"

// CurrentUser returns the current user from the session.
func CurrentUser(c fiber.Ctx) *db.User {
	user, _ := c.Locals(CurrentUserKey).(*db.User)
	return user
}

// IsCurrentUser returns whether user is the current user.
func IsCurrentUser(c fiber.Ctx, userID int64) bool {
	return CurrentUser(c).ID == userID
}

// Get returns the encapsulated globals to be referenced from templates.
func Get(c fiber.Ctx) *Values {
	values := Values{
		User:         CurrentUser(c),
		CurrentRoute: c.Route().Name,
	}

	if flashMessage, found := session.GetOnce[flash.Message](c, "globals.flash"); found {
		values.Flash = &flashMessage
	}

	return &values
}
