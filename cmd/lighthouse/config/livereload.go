package config

import "github.com/google/uuid"

var bootID = uuid.NewString()

func BootID() string {
	return bootID
}
