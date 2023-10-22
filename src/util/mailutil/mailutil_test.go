package mailutil

// This is one of the rare situations where the test code sits within the same package as the code
// it's testing. This is so that sendFunc can be tested without exporting it.

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendFunc(t *testing.T) {
	testCases := []struct {
		inputSubject string
		inputBody    string
	}{
		{"Subject", "Body"},
	}

	fromHeader := "lighthouse <lighthouse@example.com>"
	toHeader := "alerts <alerts@example.com>"

	Init(&Config{
		FromHeader: fromHeader,
		ToHeader:   toHeader,
	})

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputSubject, testCase.inputBody), func(t *testing.T) {
			mockSend := func(_ string, _ smtp.Auth, _ string, _ []string, message []byte) error {
				expectedBody := fmt.Sprintf(
					"From: %s\nTo: %s\nSubject: %s\n\n%s",
					fromHeader,
					toHeader,
					testCase.inputSubject,
					testCase.inputBody,
				)

				assert.Equal(t, expectedBody, string(message))
				return nil
			}

			err := sendFunc(mockSend, testCase.inputSubject, testCase.inputBody)
			assert.Nil(t, err)
		})
	}
}
