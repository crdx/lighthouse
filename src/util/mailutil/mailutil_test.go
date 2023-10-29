package mailutil

// This is one of the rare situations where the test code sits within the same package as the code
// it's testing. This is so that sendFunc can be tested without exporting it.

import (
	"fmt"
	"net/smtp"
	"testing"

	"crdx.org/lighthouse/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	env.Init()

	testCases := []struct {
		inputConfig *Config
		expectErr   bool
	}{
		{
			&Config{Enabled: func() bool { return false }},
			false,
		},
		{
			nil,
			true,
		},
		{
			&Config{
				SendToStdErr: false,
				Enabled:      func() bool { return true },
				Host:         func() string { return "5750bbfc-bd61-4f8f-a11f-52a345cd2d98" },
				Port:         func() string { return "1234" },
				User:         func() string { return "" },
				Pass:         func() string { return "" },
				FromAddress:  func() string { return "" },
				ToAddress:    func() string { return "" },
				FromHeader:   func() string { return "" },
				ToHeader:     func() string { return "" },
			},
			true,
		},
		{
			&Config{
				SendToStdErr: true,
				Enabled:      func() bool { return true },
				FromHeader:   func() string { return "" },
				ToHeader:     func() string { return "" },
			},
			false,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			Init(testCase.inputConfig)
			actual := Send("", "")

			if testCase.expectErr {
				require.Error(t, actual)
			} else {
				require.NoError(t, actual)
			}
		})
	}
}

func TestSendFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputSubject string
		inputBody    string
	}{
		{"Subject", "Body"},
	}

	fromHeader := "lighthouse <lighthouse@example.com>"
	toHeader := "alerts <alerts@example.com>"

	Init(&Config{
		Enabled:     func() bool { return true },
		Host:        func() string { return "smtp.example.com" },
		Port:        func() string { return "587" },
		User:        func() string { return "username" },
		Pass:        func() string { return "hunter2" },
		FromAddress: func() string { return "lighthouse@example.com" },
		ToAddress:   func() string { return "alerts@example.com" },
		FromHeader:  func() string { return fromHeader },
		ToHeader:    func() string { return toHeader },
	})

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputSubject, testCase.inputBody), func(t *testing.T) {
			t.Parallel()

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
			require.NoError(t, err)
		})
	}
}
