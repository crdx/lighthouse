package userR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

const (
	RoleNone   uint = 0
	RoleViewer uint = 1
	RoleEditor uint = 2
	RoleAdmin  uint = 3
)

func Map() map[uint]*m.User {
	users := map[uint]*m.User{}

	for _, user := range db.B[m.User]().Find() {
		users[user.ID] = user
	}

	return users
}
