package runtimeutil

// A temporary home for all methods that don't yet have a better place to go.

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/samber/lo"
)

// FindProjectRoot walks up the directory tree to find the project root (where go.mod exists).
func FindProjectRoot() string {
	dir := lo.Must(os.Getwd())

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod not found)")
		}
		dir = parent
	}
}

// PrintStackTrace prints out the current stack trace, stripping out stripDepth levels from the
// top.
func PrintStackTrace(stripDepth int) {
	b := make([]byte, 1024*10)
	length := runtime.Stack(b, false)
	s := string(b[:length])
	lines := strings.Split(s, "\n")
	header := lines[0]
	lines = lines[1+stripDepth*2:]
	fmt.Println(header)
	fmt.Println(strings.Join(lines, "\n"))
}
