package netutil_test

import (
	"fmt"
	"net"
	"testing"

	"crdx.org/lighthouse/util/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVendor(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			actualVendor, actualFound := netutil.GetVendor(testCase.inputMACAddress)
			assert.Equal(t, testCase.expectedVendor, actualVendor)
			assert.Equal(t, testCase.expectedFound, actualFound)
		})
	}
}

func TestUnqualifyHostname(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			actualHostname := netutil.UnqualifyHostname(testCase.inputHostname)
			assert.Equal(t, testCase.expectedHostname, actualHostname)
		})
	}
}

func TestExpandIPNet(t *testing.T) {
	t.Parallel()

	generate := func(prefix string, n int) []string {
		var ips []string

		for i := 1; i < n; i++ {
			ips = append(ips, prefix+fmt.Sprint(i))
		}

		return ips
	}

	testCases := []struct {
		inputIPNet  string
		inputMask   string
		expectedIPs []string
	}{
		{"192.168.1.1", "24", generate("192.168.1.", 255)},
		{"192.168.1.1", "25", generate("192.168.1.", 127)},
		{"192.168.1.1", "26", generate("192.168.1.", 63)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s/%s", testCase.inputIPNet, testCase.inputMask), func(t *testing.T) {
			t.Parallel()

			_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%s", testCase.inputIPNet, testCase.inputMask))
			require.NoError(t, err)

			actualIPs := netutil.ExpandIPNet(ipNet)
			var actualIPsStr []string
			for _, ip := range actualIPs {
				actualIPsStr = append(actualIPsStr, ip.String())
			}
			assert.Equal(t, testCase.expectedIPs, actualIPsStr)
		})
	}
}

func TestIPNetTooLarge(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputIPNet     string
		inputFixedBits int
		expected       bool
	}{
		{"192.168.0.0", 24, false},
		{"192.168.0.0", 15, true},
		{"10.0.0.0", 8, true},
		{"172.16.0.0", 16, false},
		{"172.16.0.0", 12, true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s/%d", testCase.inputIPNet, testCase.inputFixedBits), func(t *testing.T) {
			t.Parallel()

			_, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", testCase.inputIPNet, testCase.inputFixedBits))
			actual := netutil.IPNetTooLarge(ipNet)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
