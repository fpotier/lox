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
	Value lexer.Literal
}

func NewLiteralExpression(literal lexer.Literal) *LiteralExpression {
	return &LiteralExpression{
		Value: literal,
	}
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
