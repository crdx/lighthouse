package dbutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/util/dbutil"
	"github.com/stretchr/testify/assert"
)

func TestMapByID(t *testing.T) {
	type S struct {
		ID   uint
		Name string
	}

	testCases := []struct {
		items    []*S
		expected map[uint]*S
	}{
		{
			[]*S{{1, "Alice"}, {2, "Bob"}},
			map[uint]*S{1: {1, "Alice"}, 2: {2, "Bob"}},
		},
		{
			[]*S{},
			map[uint]*S{},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			actual := dbutil.MapByID(testCase.items)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMapBy(t *testing.T) {
	type S struct {
		ID   uint
		Name string
	}

	testCases := []struct {
		keyField string
		items    []*S
		expected map[uint]*S
	}{
		{
			"ID",
			[]*S{{1, "Alice"}, {2, "Bob"}},
			map[uint]*S{1: {1, "Alice"}, 2: {2, "Bob"}},
		},
		{
			"ID",
			[]*S{},
			map[uint]*S{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.keyField, func(t *testing.T) {
			actual := dbutil.MapBy[uint, S](testCase.keyField, testCase.items)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestMapBy2(t *testing.T) {
	type S struct {
		ID   uint
		Name string
	}

	testCases := []struct {
		keyField   string
		valueField string
		items      []*S
		expected   map[uint]string
	}{
		{
			"ID",
			"Name",
			[]*S{{1, "Alice"}, {2, "Bob"}},
			map[uint]string{1: "Alice", 2: "Bob"},
		},
		{
			"ID",
			"Name",
			[]*S{},
			map[uint]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.keyField, func(t *testing.T) {
			actual := dbutil.MapBy2[uint, string](testCase.keyField, testCase.valueField, testCase.items)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
