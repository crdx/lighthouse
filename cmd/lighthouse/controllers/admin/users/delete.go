package users

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v3"
)

func Delete(c fiber.Ctx) error {
	user := parseparam.Get[db.User](c)

	// Current user can't delete themselves.
	if globals.IsCurrentUser(c, user.ID) {
		return c.SendStatus(400)
	}

	user.Delete()

	auditLogR.Add(c, "Deleted user %s", user.AuditName())
	flash.Success(c, "User deleted")
	return c.Redirect().To("/admin/users")
}
