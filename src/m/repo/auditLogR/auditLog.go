package auditLogR

import (
	"fmt"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/util/sqlutil"
	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx, message string, args ...any) {
	var userID sqlutil.NullUint
	if user := globals.CurrentUser(c); user != nil {
		_ = userID.Scan(user.ID)
	}

	ipAddress := c.IP()

	var deviceID sqlutil.NullUint
	if adapter, found := db.B[m.Adapter]("ip_address = ?", ipAddress).First(); found {
		_ = deviceID.Scan(adapter.ID)
	}

	db.Save(&m.AuditLog{
		IPAddress: ipAddress,
		DeviceID:  deviceID,
		UserID:    userID,
		Message:   fmt.Sprintf(message, args...),
	})
}
