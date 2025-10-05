package device

import (
	"fmt"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"github.com/gofiber/fiber/v3"
)

func Merge(c fiber.Ctx) error {
	device1 := parseparam.Get[db.Device](c)
	device2, found := db.FindDevice(c.FormValue("device_id"))

	if !found {
		return c.SendStatus(400)
	}

	// The parent device will be the one that was discovered first.
	var parent, child *db.Device
	if device1.CreatedAt.Before(device2.CreatedAt) {
		parent, child = device1, device2
	} else {
		parent, child = device2, device1
	}

	for _, adapter := range child.Adapters() {
		adapter.UpdateDeviceID(parent.ID)
	}

	if child.LastSeenAt.After(parent.LastSeenAt) {
		parent.UpdateLastSeenAt(child.LastSeenAt)
	}

	db.MigrateDeviceStateLogs(db.MigrateDeviceStateLogsParams{
		ToDeviceID:   parent.ID,
		FromDeviceID: child.ID,
	})

	db.MigrateDeviceStateNotifications(db.MigrateDeviceStateNotificationsParams{
		ToDeviceID:   parent.ID,
		FromDeviceID: child.ID,
	})

	db.MigrateDeviceDiscoveryNotifications(db.MigrateDeviceDiscoveryNotificationsParams{
		ToDeviceID:   parent.ID,
		FromDeviceID: child.ID,
	})

	child.Delete()

	auditLogR.Add(c, "Merged device %s into %s", child.AuditName(), parent.AuditName())

	flash.Success(c,
		"Device %s merged into %s",
		child.Identifier(),
		parent.Identifier(),
	)

	return c.Redirect().To(fmt.Sprintf("/device/%d", parent.ID))
}
