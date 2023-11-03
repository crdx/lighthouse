package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func IconToClass(icon string) string {
	if style, name, ok := strings.Cut(icon, ":"); ok {
		return fmt.Sprintf("fa-%s fa-%s", style, name) //nolint
	} else {
		return ""
	}
}
