package dbutil

import "strings"

// QueryBuilder allows building up a query piecemeal while still providing the security of prepared
// statements.
type QueryBuilder struct {
	query    []string
	args     []any
	hasWhere bool
}

// NewQueryBuilder returns a new *QueryBuilder with the initial query set to query.
func NewQueryBuilder(query string) *QueryBuilder {
	return &QueryBuilder{
		query: []string{query},
	}
}

// Append adds a string to the query.
func (self *QueryBuilder) Append(s string, args ...any) {
	self.query = append(self.query, s)
	self.args = append(self.args, args...)
}

// And adds an AND condition to the query if a WHERE clause already exists, otherwise adds a WHERE
// clause.
func (self *QueryBuilder) And(condition string, args ...any) {
	self.where("AND", condition, args...)
}

// And adds an OR condition to the query if a WHERE clause already exists, otherwise adds a WHERE
// clause.
func (self *QueryBuilder) Or(condition string, args ...any) {
	self.where("OR", condition, args...)
}

// Query returns the query as a string.
func (self *QueryBuilder) Query() string {
	return strings.TrimSpace(strings.Join(self.query, " "))
}

// Args returns the number of arguments that correspond to the number of placeholders in the query.
func (self *QueryBuilder) Args() []any {
	return self.args
}

func (self *QueryBuilder) where(andor string, condition string, args ...any) {
	if !self.hasWhere {
		self.query = append(self.query, "WHERE")
		self.hasWhere = true
	} else {
		self.query = append(self.query, andor)
	}
	self.query = append(self.query, condition)
	self.args = append(self.args, args...)
}
