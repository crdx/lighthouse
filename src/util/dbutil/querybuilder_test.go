package dbutil_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/util/dbutil"
	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	testCases := []struct {
		actions       []func(*dbutil.QueryBuilder)
		expectedQuery string
		expectedArgs  []any
	}{
		{
			actions:       []func(*dbutil.QueryBuilder){},
			expectedQuery: "",
			expectedArgs:  []any{},
		},
		{
			actions: []func(*dbutil.QueryBuilder){
				func(q *dbutil.QueryBuilder) { q.Append("SELECT * FROM table") },
			},
			expectedQuery: "SELECT * FROM table",
			expectedArgs:  []any{},
		},
		{
			actions: []func(*dbutil.QueryBuilder){
				func(q *dbutil.QueryBuilder) { q.Append("SELECT * FROM table") },
				func(q *dbutil.QueryBuilder) { q.And("id = ?", 1) },
			},
			expectedQuery: "SELECT * FROM table WHERE id = ?",
			expectedArgs:  []any{1},
		},
		{
			actions: []func(*dbutil.QueryBuilder){
				func(q *dbutil.QueryBuilder) { q.Append("SELECT * FROM table") },
				func(q *dbutil.QueryBuilder) { q.And("id = ?", 1) },
				func(q *dbutil.QueryBuilder) { q.Or("name = ?", "John") },
			},
			expectedQuery: "SELECT * FROM table WHERE id = ? OR name = ?",
			expectedArgs:  []any{1, "John"},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
			q := dbutil.NewQueryBuilder("")
			for _, action := range testCase.actions {
				action(q)
			}

			assert.Equal(t, testCase.expectedQuery, q.Query())

			if len(testCase.expectedArgs) == 0 && len(q.Args()) == 0 {
				assert.True(t, true)
			} else {
				assert.Equal(t, testCase.expectedArgs, q.Args())
			}
		})
	}
}
