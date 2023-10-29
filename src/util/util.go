package util

import (
	"errors"
	"os"
)

func PathExists(str string) bool {
	_, err := os.Stat(str)
	return !errors.Is(err, os.ErrNotExist)
}
