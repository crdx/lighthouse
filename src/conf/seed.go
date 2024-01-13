package conf

import (
	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/stringutil"
)

func seed() error {
	db.Save(&m.Setting{Name: "timezone", Value: "Europe/London"})
	db.Save(&m.Setting{Name: "notify_on_new_device", Value: "1"})
	db.Save(&m.Setting{Name: "ping_new", Value: "1"})
	db.Save(&m.Setting{Name: "watch_new", Value: "0"})

	db.Save(&m.Setting{Name: "enable_device_scan", Value: "1"})
	db.Save(&m.Setting{Name: "device_scan_interval", Value: "1 min"})

	db.Save(&m.Setting{Name: "enable_service_scan", Value: "0"})
	db.Save(&m.Setting{Name: "service_scan_interval", Value: "2 hours"})

	db.Save(&m.User{Username: "root", PasswordHash: stringutil.Hash(env.DefaultRootPass()), Role: constants.RoleAdmin})
	db.Save(&m.User{Username: "anon", PasswordHash: stringutil.Hash(env.DefaultAnonPass()), Role: constants.RoleViewer})

	return nil
}
