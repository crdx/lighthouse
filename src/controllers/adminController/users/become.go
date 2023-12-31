package users

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

func Become(c *fiber.Ctx) error {
	user := util.Param[m.User](c)
	session.Set(c, "user_id", user.ID)
	flash.Success(c, "Became %s", user.Username)
	auditLogR.Add(c, "Became %s", user.Username)
	return c.Redirect("/")
}
