package query

import "time"

// Condition to marker interface.
type Condition interface {
	isCondition()
}

type Operator string

const (
	OpEq  Operator = "EQ"
	OpNeq Operator = "NEQ"
	OpGt  Operator = "GT"
	OpLt  Operator = "LT"
	OpIn  Operator = "IN"
)

type Predicate struct {
	Field    string
	Operator Operator
	Value    any
}

func (Predicate) isCondition() {}

type LogicOp string

const (
	LogicAnd LogicOp = "AND"
	LogicOr  LogicOp = "OR"
)

type Group struct {
	Operator   LogicOp
	Conditions []Condition
}

func (Group) isCondition() {}

func And(conds ...Condition) Group {
	return Group{Operator: LogicAnd, Conditions: conds}
}

func Or(conds ...Condition) Group {
	return Group{Operator: LogicOr, Conditions: conds}
}

type StringField string

func (f StringField) Eq(v string) Condition    { return Predicate{string(f), OpEq, v} }
func (f StringField) In(v ...string) Condition { return Predicate{string(f), OpIn, v} }

type IntField string

func (f IntField) Gt(v int) Condition { return Predicate{string(f), OpGt, v} }
func (f IntField) Eq(v int) Condition { return Predicate{string(f), OpEq, v} }

type TimeField string

func (f TimeField) Gt(v time.Time) Condition { return Predicate{string(f), OpGt, v} }
func (f TimeField) Eq(v time.Time) Condition { return Predicate{string(f), OpEq, v} }

