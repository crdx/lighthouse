package runtimeutil_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"crdx.org/lighthouse/pkg/util/runtimeutil"
	"github.com/stretchr/testify/assert"
)

func TestPrintStackTrace(t *testing.T) {
	testCases := []struct {
		input            int
		expectedContains []string
	}{
		{0, []string{"goroutine", "TestPrintStackTrace"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", testCase.input), func(t *testing.T) {
			originalStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			runtimeutil.PrintStackTrace(testCase.input)

			_ = w.Close()
			os.Stderr = originalStderr

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			for _, expected := range testCase.expectedContains {
				assert.Contains(t, output, expected)
			}
		})
	}
}
