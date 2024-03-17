package config

import (
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/stringutil"
)

func seed() {
	db.CreateSetting(&db.Setting{Name: "timezone", Value: "Europe/London"})
	db.CreateSetting(&db.Setting{Name: "notify_on_new_device", Value: "1"})
	db.CreateSetting(&db.Setting{Name: "ping_new", Value: "1"})
	db.CreateSetting(&db.Setting{Name: "watch_new", Value: "0"})

	db.CreateSetting(&db.Setting{Name: "enable_device_scan", Value: "1"})
	db.CreateSetting(&db.Setting{Name: "device_scan_interval", Value: "1 min"})

	db.CreateSetting(&db.Setting{Name: "enable_service_scan", Value: "0"})
	db.CreateSetting(&db.Setting{Name: "service_scan_interval", Value: "2 hours"})

	db.CreateUser(&db.User{Username: "root", PasswordHash: stringutil.Hash(env.DefaultRootPass()), Role: constants.RoleAdmin})
	db.CreateUser(&db.User{Username: "anon", PasswordHash: stringutil.Hash(env.DefaultAnonPass()), Role: constants.RoleViewer})
}
