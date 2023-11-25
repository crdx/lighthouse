package users

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	user := util.Param[m.User](c)

	// Current user can't delete themselves.
	if globals.IsCurrentUser(c, user) {
		return c.SendStatus(400)
	}

	user.Delete()

	auditLogR.Add(c, "Deleted user %s", user.AuditName())
	flash.Success(c, "User deleted")
	return c.Redirect("/admin/users")
}
