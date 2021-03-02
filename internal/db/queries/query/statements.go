package query

import "fmt"

// Table defines a table location for a given query.
type Table string

// Where clause for a given query.
type Where Expression

// Statements defines generic queries for a given underlying query engine.
type Statements interface {
	Builder

	// Count returns the number of rows in the given table.
	Count(Table, Where) Query

	// Delete removes the row identified by the given ID.
	Delete(Table) Query

	// Upsert inserts or replaces a new row with the given column values,
	// to the given table using columns order.
	Upsert(Table, []string) Query
}

// Equals defines an equality expression for the where clause.
func Equals(name string, value interface{}) Expression {
	return Statement{
		Name:  name,
		Value: value,
		Op:    Equal,
	}
}

// And defines a AND expression for the where clause.
func And(right, left Expression) Expression {
	return Infix{
		Right: right,
		Left:  left,
		Op:    AND,
	}
}

// ExpressionOperatorType defines a join operator type for constructing
// multiple expressions in a where clause.
type ExpressionOperatorType string

const (
	// AND defines a AND operator.
	AND ExpressionOperatorType = "AND"
	// OR defines a OR operator.
	OR ExpressionOperatorType = "OR"
	// NOT defines a NOT operator.
	NOT ExpressionOperatorType = "NOT"
)

// OperatorType defines a operator type for constructing where clauses.
type OperatorType string

const (
	// Equal defines a = operator.
	Equal OperatorType = "="
)

// Builder defines an interface for creating the correct params and operators
// dependant on the underlying interface.
type Builder interface {
	Params(int) string
	Op(OperatorType) string
	ExpressionOp(ExpressionOperatorType) string
}

// Expression defines a way of constructing a where clause.
type Expression interface {
	Build(Builder) (string, []interface{})
}

type Statement struct {
	Name  string
	Value interface{}
	Op    OperatorType
}

func (e Statement) Build(builder Builder) (string, []interface{}) {
	return fmt.Sprintf("%s %s %s", e.Name, builder.Op(e.Op), builder.Params(1)), []interface{}{e.Value}
}

type Infix struct {
	Right Expression
	Left  Expression
	Op    ExpressionOperatorType
}

func (s Infix) Build(builder Builder) (string, []interface{}) {
	rs, rp := s.Right.Build(builder)
	ls, lp := s.Left.Build(builder)
	return fmt.Sprintf("%s %s %s", rs, builder.ExpressionOp(s.Op), ls), append(rp, lp)
}
