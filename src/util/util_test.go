package util_test

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
	"testing"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/util"
	"github.com/stretchr/testify/assert"
)

func TestPrintStackTrace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input            int
		expectedContains []string
	}{
		{0, []string{"goroutine", "TestPrintStackTrace"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", testCase.input), func(t *testing.T) {
			t.Parallel()

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

func TestPluralise(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputCount int
		inputUnit  string
		expected   string
	}{
		{1, "apple", "apple"},
		{2, "apple", "apples"},
		{0, "apple", "apples"},
		{1, "item", "item"},
		{2, "item", "items"},
		{-1, "apple", "apples"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%s", testCase.inputCount, testCase.inputUnit), func(t *testing.T) {
			t.Parallel()

			actual := util.Pluralise(testCase.inputCount, testCase.inputUnit)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestSendMail(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputSubject string
		inputMessage string
	}{
		{"Subject", "Message"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputSubject, testCase.inputMessage), func(t *testing.T) {
			t.Parallel()

			mockSend := func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
				expectedBody := fmt.Sprintf(
					"To: %s\nSubject: %s\n\n%s",
					env.MailTo,
					testCase.inputSubject,
					testCase.inputMessage,
				)

				assert.Equal(t, expectedBody, string(msg))
				return nil
			}

			err := util.SendMailFunc(mockSend, testCase.inputSubject, testCase.inputMessage)
			assert.Nil(t, err)
		})
	}
}
