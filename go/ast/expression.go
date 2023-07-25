package ast

import "github.com/fpotier/crafting-interpreters/go/lexer"

type Expression interface {
	Accept(visitor Visitor)
}

type BinaryExpression struct {
	LHS      Expression
	Operator lexer.Token
	RHS      Expression
}

func NewBinaryExpression(lhs Expression, operator lexer.Token, rhs Expression) *BinaryExpression {
	return &BinaryExpression{
		LHS:      lhs,
		Operator: operator,
		RHS:      rhs,
	}
}

func (e *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(e)
}

type GroupingExpression struct {
	Expr Expression
}

func NewGroupingExpression(expr Expression) *GroupingExpression {
	return &GroupingExpression{
		Expr: expr,
	}
}

func (e *GroupingExpression) Accept(visitor Visitor) {
	visitor.VisitGroupingExpression(e)
}

type LiteralExpression struct {
	// TODO this should be a LoxValue of something like that
	value LoxValue
}

func NewLiteralExpression(value LoxValue) *LiteralExpression {
	return &LiteralExpression{
		value: value,
	}
}

func (e *LiteralExpression) LoxValue() LoxValue {
	return e.value
}

func (e *LiteralExpression) Accept(visitor Visitor) {
	visitor.VisitLiteralExpression(e)
}

type UnaryExpression struct {
	Operator lexer.Token
	RHS      Expression
}

func NewUnaryExpression(operator lexer.Token, rhs Expression) *UnaryExpression {
	return &UnaryExpression{
		Operator: operator,
		RHS:      rhs,
	}
}

func (e *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(e)
}

type VariableExpression struct {
	Name lexer.Token
}

func NewVariableExpression(name lexer.Token) *VariableExpression {
	return &VariableExpression{Name: name}
}

func (e *VariableExpression) Accept(visitor Visitor) {
	visitor.VisitVariableExpression(e)
}

type AssignmentExpression struct {
	Name  lexer.Token
	Value Expression
}

func NewAssignmentExpression(name lexer.Token, value Expression) *AssignmentExpression {
	return &AssignmentExpression{
		Name:  name,
		Value: value,
	}
}

func (e *AssignmentExpression) Accept(visitor Visitor) {
	visitor.VisitAssignmentExpression(e)
}

type LogicalExpression struct {
	LHS      Expression
	Operator lexer.Token
	RHS      Expression
}

func NewLogicalExpression(lhs Expression, operator lexer.Token, rhs Expression) *LogicalExpression {
	return &LogicalExpression{
		LHS:      lhs,
		Operator: operator,
		RHS:      rhs,
	}
}

func (e *LogicalExpression) Accept(visitor Visitor) {
	visitor.VisitLogicalExpression(e)
}

type CallExpression struct {
	Callee   Expression
	Position lexer.Token
	Args     []Expression
}

func NewCallExpression(callee Expression, position lexer.Token, args []Expression) *CallExpression {
	return &CallExpression{
		Callee:   callee,
		Position: position,
		Args:     args,
	}
}

func (e *CallExpression) Accept(visitor Visitor) {
	visitor.VisitCallExpression(e)
}
