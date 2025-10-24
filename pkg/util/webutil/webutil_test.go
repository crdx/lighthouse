package webutil_test

import (
	"testing"

	"crdx.org/lighthouse/pkg/util/webutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		basePath    string
		queryParams map[string]string
		expected    string
		expectError bool
	}{
		{"https://example.com", map[string]string{"key1": "value1"}, "https://example.com?key1=value1", false},
		{"https://example.com", map[string]string{"key1": "value1", "key2": "value2"}, "https://example.com?key1=value1&key2=value2", false},
		{"https://example.com", map[string]string{}, "https://example.com", false},
		{"https://example.com", nil, "https://example.com", false},
		{":", map[string]string{"key1": "value1"}, "", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.basePath, func(t *testing.T) {
			t.Parallel()

			actual, err := webutil.BuildURL(testCase.basePath, testCase.queryParams)

			if testCase.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
