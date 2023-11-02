package util

import (
	"errors"
	"os"
)

func PathExists(str string) bool {
	_, err := os.Stat(str)
	return !errors.Is(err, os.ErrNotExist)
}

// Chain runs the provided functions until it reaches one that returns a non-nil error, then returns
// it. Returns nil if none of the functions errored.
func Chain(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
