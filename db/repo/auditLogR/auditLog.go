package auditLogR

import (
	"database/sql"
	"fmt"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/globals"
	"github.com/gofiber/fiber/v3"
)

func Add(c fiber.Ctx, message string, args ...any) {
	var userID sql.Null[int64]
	if user := globals.CurrentUser(c); user != nil {
		userID = db.N(user.ID)
	}

	ipAddress := c.IP()

	var deviceID sql.Null[int64]
	if adapter, found := db.FindAdapterByIPAddressLatest(ipAddress); found {
		deviceID = db.N(adapter.DeviceID)
	}

	db.CreateAuditLog(&db.AuditLog{
		IPAddress: ipAddress,
		DeviceID:  deviceID,
		UserID:    userID,
		Message:   fmt.Sprintf(message, args...),
	})
}
