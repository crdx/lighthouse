package util_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"crdx.org/lighthouse/util"
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
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			util.PrintStackTrace(testCase.input)

			w.Close()
			os.Stdout = originalStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			for _, expected := range testCase.expectedContains {
				assert.True(t, strings.Contains(output, expected))
			}
		})
	}
}
