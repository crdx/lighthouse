package tplutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/tplutil"
	"github.com/go-playground/assert/v2"
)

func TestAddSortMetadata(t *testing.T) {
	testCases := []struct {
		currentSortColumn    string
		currentSortDirection string
		input                map[string]tplutil.SortableColumnConfig
		expected             map[string]tplutil.SortableColumnState
	}{
		{
			"column1",
			"asc",
			map[string]tplutil.SortableColumnConfig{
				"column1": {Label: "Column 1", DefaultSortDirection: "asc", Minimal: false},
				"column2": {Label: "Column 2", DefaultSortDirection: "asc", Minimal: false},
			},
			map[string]tplutil.SortableColumnState{
				"column1": {Label: "Column 1", CurrentSortColumn: "column1", CurrentSortDirection: "asc", SortColumn: "column1", SortDirection: "desc", Minimal: false},
				"column2": {Label: "Column 2", CurrentSortColumn: "column1", CurrentSortDirection: "asc", SortColumn: "column2", SortDirection: "asc", Minimal: false},
			},
		},
		{
			"column2",
			"desc",
			map[string]tplutil.SortableColumnConfig{
				"column1": {Label: "Column 1", DefaultSortDirection: "asc", Minimal: true},
				"column2": {Label: "Column 2", DefaultSortDirection: "desc", Minimal: true},
			},
			map[string]tplutil.SortableColumnState{
				"column1": {Label: "Column 1", CurrentSortColumn: "column2", CurrentSortDirection: "desc", SortColumn: "column1", SortDirection: "asc", Minimal: true},
				"column2": {Label: "Column 2", CurrentSortColumn: "column2", CurrentSortDirection: "desc", SortColumn: "column2", SortDirection: "asc", Minimal: true},
			},
		},
		{
			"",
			"",
			map[string]tplutil.SortableColumnConfig{
				"column1": {Label: "Column 1", DefaultSortDirection: "asc", Minimal: false},
			},
			map[string]tplutil.SortableColumnState{
				"column1": {Label: "Column 1", CurrentSortColumn: "", CurrentSortDirection: "", SortColumn: "column1", SortDirection: "asc", Minimal: false},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			actual := tplutil.AddSortMetadata(
				testCase.currentSortColumn,
				testCase.currentSortDirection,
				testCase.input,
			)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
