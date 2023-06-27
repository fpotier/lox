package ast

import "github.com/fpotier/crafting-interpreters/go/lexer"

type Expression interface {
	Accept(visitor Visitor)
}

type BinaryExpression struct {
	Lhs      Expression
	Operator lexer.Token
	Rhs      Expression
}

func NewBinaryExpression(lhs Expression, operator lexer.Token, rhs Expression) *BinaryExpression {
	return &BinaryExpression{
		Lhs:      lhs,
		Operator: operator,
		Rhs:      rhs,
	}
}

func (expr *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(expr)
}

type GroupingExpression struct {
	Expr Expression
}

func NewGroupingExpression(expr Expression) *GroupingExpression {
	return &GroupingExpression{
		Expr: expr,
	}
}

func (expr *GroupingExpression) Accept(visitor Visitor) {
	visitor.VisitGroupingExpression(expr)
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

func (literalExpression *LiteralExpression) LoxValue() LoxValue {
	return literalExpression.value
}

func (expr *LiteralExpression) Accept(visitor Visitor) {
	visitor.VisitLiteralExpression(expr)
}

type UnaryExpression struct {
	Operator lexer.Token
	Rhs      Expression
}

func NewUnaryExpression(operator lexer.Token, rhs Expression) *UnaryExpression {
	return &UnaryExpression{
		Operator: operator,
		Rhs:      rhs,
	}
}

func (expr *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(expr)
}

type VariableExpression struct {
	Name lexer.Token
}

func NewVariableExpression(name lexer.Token) *VariableExpression {
	return &VariableExpression{Name: name}
}

func (expr *VariableExpression) Accept(visitor Visitor) {
	visitor.VisitVariableExpression(expr)
}
