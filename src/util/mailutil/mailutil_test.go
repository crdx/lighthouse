package mailutil

// This is one of the rare situations where the test code sits within the same package as the code
// it's testing. This is so that sendFunc can be tested without exporting it.

import (
	"fmt"
	"net/smtp"
	"testing"

	"crdx.org/lighthouse/setting"
	"github.com/stretchr/testify/assert"
)

func TestSendFunc(t *testing.T) {
	testCases := []struct {
		inputSubject string
		inputBody    string
	}{
		{"Subject", "Body"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputSubject, testCase.inputBody), func(t *testing.T) {
			mockSend := func(_ string, _ smtp.Auth, _ string, _ []string, message []byte) error {
				expectedBody := fmt.Sprintf(
					"From: %s\nTo: %s\nSubject: %s\n\n%s",
					setting.Get(setting.NotificationFromHeader),
					setting.Get(setting.NotificationToHeader),
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
