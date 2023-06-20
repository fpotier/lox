package ast

import (
	"github.com/fpotier/crafting-interpreters/go/lexer"
)

type Expression interface {
	Accept(visitor Visitor) interface{}
}

type BinaryExpression struct {
	lhs      Expression
	operator lexer.Token
	rhs      Expression
}

func NewBinaryExpression(lhs Expression, operator lexer.Token, rhs Expression) *BinaryExpression {
	return &BinaryExpression{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

func (expr *BinaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpression(expr)
}

type GroupingExpression struct {
	expr Expression
}

func NewGroupingExpression(expr Expression) *GroupingExpression {
	return &GroupingExpression{
		expr: expr,
	}
}

func (expr *GroupingExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpression(expr)
}

type LiteralExpression struct {
	value lexer.Literal
}

func NewLiteralExpression(literal lexer.Literal) *LiteralExpression {
	return &LiteralExpression{
		value: literal,
	}
}

func (expr *LiteralExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpression(expr)
}

type UnaryExpression struct {
	operator lexer.Token
	rhs      Expression
}

func NewUnaryExpression(operator lexer.Token, rhs Expression) *UnaryExpression {
	return &UnaryExpression{
		operator: operator,
		rhs:      rhs,
	}
}

func (expr *UnaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpression(expr)
}
