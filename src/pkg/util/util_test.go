package util_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIconToClass(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"solid:home", "fa-solid fa-home"},
		{"regular:bell", "fa-regular fa-bell"},
		{"brands:github", "fa-brands fa-github"},
		{"light:camera", "fa-light fa-camera"},
		{"", ""},
		{"invalid", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := util.IconToClass(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestChain(t *testing.T) {
	testCases := []struct {
		name     string
		input    []func() error
		expected error
	}{
		{
			"AllSuccess",
			[]func() error{
				func() error { return nil },
				func() error { return nil },
				func() error { return nil },
			},
			nil,
		},
		{
			"FirstFails",
			[]func() error{
				func() error { return fmt.Errorf("first failure") },
				func() error { return nil },
				func() error { return nil },
			},
			fmt.Errorf("first failure"),
		},
		{
			"MiddleFails",
			[]func() error{
				func() error { return nil },
				func() error { return fmt.Errorf("middle failure") },
				func() error { return nil },
			},
			fmt.Errorf("middle failure"),
		},
		{
			"LastFails",
			[]func() error{
				func() error { return nil },
				func() error { return nil },
				func() error { return fmt.Errorf("last failure") },
			},
			fmt.Errorf("last failure"),
		},
		{
			"EmptyFunctions",
			[]func() error{},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := util.Chain(testCase.input...)
			if testCase.expected == nil {
				require.NoError(t, actual)
			} else {
				require.Error(t, actual)
				assert.Equal(t, testCase.expected.Error(), actual.Error())
			}
		})
	}
}
