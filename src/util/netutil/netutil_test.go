package netutil_test

import (
	"testing"

	"crdx.org/lighthouse/util/netutil"
	"github.com/stretchr/testify/assert"
)

func TestGetVendor(t *testing.T) {
	testCases := []struct {
		inputMACAddress string
		expectedVendor  string
		expectedFound   bool
	}{
		{"00:00:00:00:00:00", "XEROX CORPORATION", true},
		{"00:1A:11:00:00:00", "Google, Inc.", true},
		{"FC:FC:48:00:00:00", "Apple, Inc.", true},
		{"12:34:56:78:9A:BC", "", false},
		{"invalid", "", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputMACAddress, func(t *testing.T) {
			actualVendor, actualFound := netutil.GetVendor(testCase.inputMACAddress)
			assert.Equal(t, testCase.expectedVendor, actualVendor)
			assert.Equal(t, testCase.expectedFound, actualFound)
		})
	}
}

func TestUnqualifyHostname(t *testing.T) {
	testCases := []struct {
		inputHostname    string
		expectedHostname string
	}{
		{"test.local.", "test"},
		{"test.local", "test"},
		{"test", "test"},
		{"test.", "test"},
		{"", ""},
		{".", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputHostname, func(t *testing.T) {
			actualHostname := netutil.UnqualifyHostname(testCase.inputHostname)
			assert.Equal(t, testCase.expectedHostname, actualHostname)
		})
	}
}
