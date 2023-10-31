package helpers

import (
	"crdx.org/db"
	"crdx.org/lighthouse/conf"
	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/middleware/auth"
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
