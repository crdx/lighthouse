package auditLogR

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx, message string, args ...any) {
	var userID uint
	if user := globals.CurrentUser(c); user != nil {
		userID = user.ID
	}

	db.Save(&m.AuditLog{
		UserID:    userID,
		IPAddress: c.IP(),
		Message:   fmt.Sprintf(message, args...),
	})
}
