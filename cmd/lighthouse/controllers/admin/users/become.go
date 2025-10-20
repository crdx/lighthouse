package users

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"crdx.org/lighthouse/pkg/session"
	"github.com/gofiber/fiber/v3"
)

func Become(c fiber.Ctx) error {
	user := parseparam.Get[db.User](c)
	session.Set(c, "user_id", user.ID)
	flash.Success(c, "Became %s", user.Username)
	auditLogR.Add(c, "Became %s", user.Username)
	return c.Redirect().To("/")
}
