package seeder

import (
	"time"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/deviceR"
	"crdx.org/lighthouse/pkg/constants"
)

func createDevice(id int64, name string, lastSeen time.Time, origin bool, state string) *db.Device {
	return db.CreateDevice(&db.Device{
		ID:          id,
		Name:        name,
		State:       state,
		Icon:        constants.DefaultDeviceIconClass,
		LastSeenAt:  lastSeen,
		CreatedAt:   lastSeen,
		Origin:      origin,
		GracePeriod: constants.DefaultGracePeriod,
	})
}

func createAdapter(id, deviceID int64, name, vendor, macAddress, ipAddress string, lastSeen time.Time) {
	db.CreateAdapter(&db.Adapter{
		ID:         id,
		DeviceID:   deviceID,
		Name:       name,
		Vendor:     vendor,
		MACAddress: macAddress,
		IPAddress:  ipAddress,
		LastSeenAt: lastSeen,
	})
}

func createDeviceStateLog(id, deviceID int64, state string, createdAt time.Time) {
	db.CreateDeviceStateLog(&db.DeviceStateLog{
		ID:          id,
		DeviceID:    deviceID,
		State:       state,
		CreatedAt:   createdAt,
		GracePeriod: "5 mins",
	})
}

func Run() {
	t1 := time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, time.October, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, time.November, 1, 12, 0, 0, 0, time.UTC)
	t4 := db.Now().Add(-1 * time.Minute)

	device1 := createDevice(1, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041", t1, false, deviceR.StateOnline)
	createAdapter(1, device1.ID, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25", "Vendor 1", "AA:AA:AA:AA:AA:AA", "127.0.0.1", t1)

	device2 := createDevice(2, "device2-64774746-5937-412c-9aa4-f262d990cc7d", t2, false, deviceR.StateOnline)
	createAdapter(2, device2.ID, "adapter2-c71739fd-d6f2-44e8-966f-fc5cdf2eec59", "Vendor 2", "BB:BB:BB:BB:BB:BB", "127.0.0.2", t2)

	device3 := createDevice(3, "device3-5acf7b73-b02c-4fe5-a63e-869f8bfc329e", t3, true, deviceR.StateOnline)
	createAdapter(3, device3.ID, "adapter3-5b083c73-f92b-4890-811a-eed7bdca99c6", "Vendor 3", "CC:CC:CC:CC:CC:CC", "127.0.0.3", t3)

	device4 := createDevice(4, "device3-8410c6ed-9b02-4c09-b757-f818f7a33586", t4, false, deviceR.StateOffline)
	createAdapter(4, device4.ID, "adapter3-3171430d-c0a0-4a11-b980-05b4f20cf699", "Vendor 4", "DD:DD:DD:DD:DD:DD", "127.0.0.4", t4)

	createDeviceStateLog(1, device1.ID, deviceR.StateOnline, time.Now().Add(-3*time.Minute))
	createDeviceStateLog(2, device1.ID, deviceR.StateOffline, time.Now().Add(-2*time.Minute))
	createDeviceStateLog(3, device2.ID, deviceR.StateOffline, time.Now().Add(-1*time.Minute))

	db.CreateNotification(&db.Notification{
		Subject: "subject-8f3fdfea-f39c-427f-b8f5-0155119975ff",
		Body:    "body-be67c77f-595c-4c91-8d24-99c829de1bbe",
	})

	// bcrypt is expensive, and this runs for every test, so pregenerate the hashes.
	rootHash := `$2a$10$Mjxj19.2lGTooLqwxi6MQeCukr7lZFyODKSFCIRR2aldNg/oTov.K`
	anonHash := `$2a$10$mnYikOcNhl.Kr4bzShVIne4vywF9zRw967qOBQpaGpbTl2HRBoCPm`
	edHash := `$2a$12$TpmujKynuODqYqFA89./iuZU.DDz1y7/K3R096NdxAgOl6fOjD9Ai`

	db.CreateUser(&db.User{ID: 1, Username: "root", PasswordHash: rootHash, Role: constants.RoleAdmin})
	db.CreateUser(&db.User{ID: 2, Username: "ed", PasswordHash: edHash, Role: constants.RoleEditor})
	db.CreateUser(&db.User{ID: 3, Username: "anon", PasswordHash: anonHash, Role: constants.RoleViewer})

	db.CreateAuditLog(&db.AuditLog{
		ID:        1,
		IPAddress: "127.0.0.1",
		UserID:    db.N[int64](1),
		Message:   "Edited device device1-625a5fa0-9b63-46d8-b4fa-578f92dca041",
	})

	db.CreateService(&db.Service{
		DeviceID: 1,
		Name:     "service1-f6d0b172-7e23-4d6c-a9bd-e456208c01fe",
		Port:     80,
	})

	db.CreateSetting(&db.Setting{
		Name:  "enable_service_scan",
		Value: "1",
	})
}
