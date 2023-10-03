package scanner

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			assert.NoError(t, err)

			actualIPs := expandIPNet(ipNet)
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
			actual := ipNetTooLarge(ipNet)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
