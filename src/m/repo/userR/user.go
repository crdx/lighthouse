package userR

import (
	"crdx.org/db"
	"crdx.org/lighthouse/m"
)

func Map() map[uint]*m.User {
	users := map[uint]*m.User{}

	for _, user := range db.B[m.User]().Find() {
		users[user.ID] = user
	}

	return users
}
