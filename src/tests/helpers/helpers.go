package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/util/mailutil"
	"crdx.org/lighthouse/util/timeutil"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// Init initialises the database and r etuns a new session with the requested auth state.
func Init(state auth.State, handlers ...func(c *fiber.Ctx) error) *Session {
	env.Init()

	dbConfig := conf.GetTestDbConfig()
	lo.Must0(db.Init(dbConfig))
	session.Init(conf.GetTestSessionConfig(), dbConfig)

	timeutil.Init(&timeutil.Config{Timezone: func() string { return "Europe/London" }})
	mailutil.Init(&mailutil.Config{Enabled: func() bool { return false }})

	return NewSession(state, handlers...)
}

// AutoLogin returns middleware that simulates the user being authorised as the provided state. The
// first user in the db with the required authorisation will be picked. This is designed to be used
// for tests.
func AutoLogin(state auth.State) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := db.B[m.User]("admin = ?", state == auth.StateAdmin).First()
		c.Locals(globals.CurrentUserKey, user)
		return c.Next()
	}
}
