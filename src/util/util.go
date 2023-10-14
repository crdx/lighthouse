package util

// A temporary home for all methods that don't yet have a better place to go.

import (
	"fmt"
	"runtime"
	"strings"
)

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
