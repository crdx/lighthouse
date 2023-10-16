package pager_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/pager"
	"github.com/stretchr/testify/assert"
)

func TestGetState(t *testing.T) {
	testCases := []struct {
		name           string
		currentPage    uint
		totalPages     uint
		path           string
		qs             map[string]string
		expectedState  pager.State
		expectedCount  uint
		expectedOffset uint
	}{
		{
			name:          "Case1",
			currentPage:   1,
			totalPages:    5,
			path:          "/example",
			qs:            map[string]string{"key": "value"},
			expectedState: pager.State{CurrentPage: 1, TotalPages: 5, NextPageURL: "/example?key=value&page=2"},
		},
		{
			name:          "Case2",
			currentPage:   3,
			totalPages:    3,
			path:          "/example",
			qs:            map[string]string{"key": "value"},
			expectedState: pager.State{CurrentPage: 3, TotalPages: 3, PreviousPageURL: "/example?key=value&page=2"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			state, err := pager.GetState(testCase.currentPage, testCase.totalPages, testCase.path, testCase.qs)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedState, *state, "state")
		})
	}
}

func TestGetPageCount(t *testing.T) {
	testCases := []struct {
		total    uint
		perPage  uint
		expected uint
	}{
		{0, 10, 0},
		{1, 10, 1},
		{9, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{100, 10, 10},
		{101, 10, 11},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%d", testCase.total, testCase.perPage), func(t *testing.T) {
			actual := pager.GetPageCount(testCase.total, testCase.perPage)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetOffset(t *testing.T) {
	testCases := []struct {
		pageNumber uint
		perPage    uint
		expected   uint
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 10, 20},
		{1, 20, 0},
		{2, 20, 20},
		{3, 5, 10},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d,%d", testCase.pageNumber, testCase.perPage), func(t *testing.T) {
			actual := pager.GetOffset(testCase.pageNumber, testCase.perPage)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}