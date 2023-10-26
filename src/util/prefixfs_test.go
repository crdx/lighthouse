package util_test

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"testing"

	"crdx.org/lighthouse/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefixFS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputPrefix string
		inputName   string
		expected    string
		expectError bool
	}{
		{"/prefix", "/file", "/prefix/file", false},
		{"/prefix", "file", "/prefix/file", false},
		{"prefix", "/file", "prefix/file", false},
		{"/prefix", "../file", "/prefix/../file", true},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%s", testCase.inputPrefix, testCase.inputName), func(t *testing.T) {
			t.Parallel()

			prefixFS := &util.PrefixFS{
				FS:     &mockFS{files: []string{"prefix/file"}},
				Prefix: testCase.inputPrefix,
			}

			file, err := prefixFS.Open(testCase.inputName)

			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, file.(*mockFile).Name)
			}
		})
	}
}

type mockFS struct {
	files []string
}

func (m *mockFS) Open(name string) (fs.File, error) {
	if !slices.Contains(m.files, strings.TrimPrefix(name, "/")) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: errors.New("file does not exist")}
	}

	return &mockFile{Name: name}, nil
}

type mockFile struct {
	Name string
}

func (*mockFile) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (*mockFile) Read(_ []byte) (n int, err error) {
	return 0, nil
}

func (*mockFile) Close() error {
	return nil
}
