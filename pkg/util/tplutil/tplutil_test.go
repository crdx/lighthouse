package tplutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/pkg/util/tplutil"
	"github.com/stretchr/testify/assert"
)

func TestAddSortMetadata(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		currentSortColumn    string
		currentSortDirection string
		queryParams          map[string]string
		input                map[string]tplutil.ColumnConfig
		expected             map[string]tplutil.ColumnState
	}{
		{
			"column1",
			"asc",
			nil,
			map[string]tplutil.ColumnConfig{
				"column1": {
					Label:                "Column 1",
					DefaultSortDirection: "asc",
					Minimal:              false,
				},
				"column2": {
					Label:                "Column 2",
					DefaultSortDirection: "asc",
					Minimal:              false,
				},
			},
			map[string]tplutil.ColumnState{
				"column1": {
					Label:                "Column 1",
					CurrentSortColumn:    "column1",
					CurrentSortDirection: "asc",
					SortColumn:           "column1",
					SortDirection:        "desc",
					Minimal:              false,
				},
				"column2": {
					Label:                "Column 2",
					CurrentSortColumn:    "column1",
					CurrentSortDirection: "asc",
					SortColumn:           "column2",
					SortDirection:        "asc",
					Minimal:              false,
				},
			},
		},
		{
			"column2",
			"desc",
			map[string]string{"f": "online"},
			map[string]tplutil.ColumnConfig{
				"column1": {
					Label:                "Column 1",
					DefaultSortDirection: "asc",
					Minimal:              true,
				},
				"column2": {
					Label:                "Column 2",
					DefaultSortDirection: "desc",
					Minimal:              true,
				},
			},
			map[string]tplutil.ColumnState{
				"column1": {
					Label:                "Column 1",
					CurrentSortColumn:    "column2",
					CurrentSortDirection: "desc",
					QueryParams:          map[string]string{"f": "online"},
					SortColumn:           "column1",
					SortDirection:        "asc",
					Minimal:              true,
				},
				"column2": {
					Label:                "Column 2",
					CurrentSortColumn:    "column2",
					CurrentSortDirection: "desc",
					QueryParams:          map[string]string{"f": "online"},
					SortColumn:           "column2",
					SortDirection:        "asc",
					Minimal:              true,
				},
			},
		},
		{
			"",
			"",
			nil,
			map[string]tplutil.ColumnConfig{
				"column1": {
					Label:                "Column 1",
					DefaultSortDirection: "asc",
					Minimal:              false,
				},
			},
			map[string]tplutil.ColumnState{
				"column1": {
					Label:                "Column 1",
					CurrentSortColumn:    "",
					CurrentSortDirection: "",
					SortColumn:           "column1",
					SortDirection:        "asc",
					Minimal:              false,
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			t.Parallel()

			actual := tplutil.AddMetadata(
				testCase.currentSortColumn,
				testCase.currentSortDirection,
				testCase.queryParams,
				testCase.input,
			)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
