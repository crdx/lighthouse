package mailutil_test

import (
	"fmt"
	"net/smtp"
	"testing"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/util/mailutil"
	"github.com/stretchr/testify/assert"
)

func TestSendNotification(t *testing.T) {
	testCases := []struct {
		inputSubject string
		inputBody    string
	}{
		{"Subject", "Body"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputSubject, testCase.inputBody), func(t *testing.T) {
			mockSend := func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
				expectedBody := fmt.Sprintf(
					"From: lighthouse (dev) <%s>\nTo: %s\nSubject: %s\n\n%s",
					env.NotificationFromAddress,
					env.NotificationToAddress,
					testCase.inputSubject,
					testCase.inputBody,
				)

				assert.Equal(t, expectedBody, string(msg))
				return nil
			}

			err := mailutil.SendNotificationFunc(mockSend, testCase.inputSubject, testCase.inputBody)
			assert.Nil(t, err)
		})
	}
}
