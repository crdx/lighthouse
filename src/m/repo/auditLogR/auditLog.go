package auditLogR

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx, message string, args ...any) {
	db.Save(&m.AuditLog{
		UserID:  globals.CurrentUser(c).ID,
		Message: fmt.Sprintf(message, args...),
	})
}
