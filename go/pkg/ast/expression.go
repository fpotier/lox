package ast

import "github.com/fpotier/lox/go/pkg/lexer"

func NewAssignmentExpression(name lexer.Token, value Expression) *AssignmentExpression {
	return &AssignmentExpression{
		Name:  name,
		Value: value,
	}
}

func NewBinaryExpression(lhs Expression, operator lexer.Token, rhs Expression) *BinaryExpression {
	return &BinaryExpression{
		LHS:      lhs,
		Operator: operator,
		RHS:      rhs,
	}
}

func NewCallExpression(callee Expression, position lexer.Token, args []Expression) *CallExpression {
	return &CallExpression{
		Callee:   callee,
		Position: position,
		Args:     args,
	}
}

func NewGetExpression(object Expression, name lexer.Token) *GetExpression {
	return &GetExpression{Object: object, Name: name}
}

func NewGroupingExpression(expr Expression) *GroupingExpression {
	return &GroupingExpression{
		Expr: expr,
	}
}

func NewLiteralExpression(value LoxValue) *LiteralExpression {
	return &LiteralExpression{
		value: value,
	}
}

func NewLogicalExpression(lhs Expression, operator lexer.Token, rhs Expression) *LogicalExpression {
	return &LogicalExpression{
		LHS:      lhs,
		Operator: operator,
		RHS:      rhs,
	}
}

func NewSetExpression(object Expression, name lexer.Token, value Expression) *SetExpression {
	return &SetExpression{
		Object: object,
		Name:   name,
		Value:  value,
	}
}

func NewSuperExpression(keyword lexer.Token, method lexer.Token) *SuperExpression {
	return &SuperExpression{Keyword: keyword, Method: method}
}

func NewThisExpression(keyword lexer.Token) *ThisExpression {
	return &ThisExpression{Keyword: keyword}
}

func NewUnaryExpression(operator lexer.Token, rhs Expression) *UnaryExpression {
	return &UnaryExpression{
		Operator: operator,
		RHS:      rhs,
	}
}

func NewVariableExpression(name lexer.Token) *VariableExpression {
	return &VariableExpression{Name: name}
}

type AssignmentExpression struct {
	Name  lexer.Token
	Value Expression
}

func (e *AssignmentExpression) Accept(visitor Visitor) { visitor.VisitAssignmentExpression(e) }

type BinaryExpression struct {
	LHS      Expression
	Operator lexer.Token
	RHS      Expression
}

func (e *BinaryExpression) Accept(visitor Visitor) { visitor.VisitBinaryExpression(e) }

type CallExpression struct {
	Callee   Expression
	Position lexer.Token
	Args     []Expression
}

func (e *CallExpression) Accept(visitor Visitor) { visitor.VisitCallExpression(e) }

type GetExpression struct {
	Object Expression
	Name   lexer.Token
}

func (e *GetExpression) Accept(visitor Visitor) { visitor.VisitGetExpression(e) }

type GroupingExpression struct {
	Expr Expression
}

func (e *GroupingExpression) Accept(visitor Visitor) { visitor.VisitGroupingExpression(e) }

type LiteralExpression struct{ value LoxValue }

func (e *LiteralExpression) LoxValue() LoxValue     { return e.value }
func (e *LiteralExpression) Accept(visitor Visitor) { visitor.VisitLiteralExpression(e) }

type LogicalExpression struct {
	LHS      Expression
	Operator lexer.Token
	RHS      Expression
}

func (e *LogicalExpression) Accept(visitor Visitor) { visitor.VisitLogicalExpression(e) }

type SetExpression struct {
	Object Expression
	Name   lexer.Token
	Value  Expression
}

func (e *SetExpression) Accept(visitor Visitor) { visitor.VisitSetExpression(e) }

type ThisExpression struct {
	Keyword lexer.Token
}

func (e *ThisExpression) Accept(visitor Visitor) { visitor.VisitThisExpression(e) }

type SuperExpression struct {
	Keyword lexer.Token
	Method  lexer.Token
}

func (e *SuperExpression) Accept(visitor Visitor) { visitor.VisitSuperExpression(e) }

type UnaryExpression struct {
	Operator lexer.Token
	RHS      Expression
}

func (e *UnaryExpression) Accept(visitor Visitor) { visitor.VisitUnaryExpression(e) }

type VariableExpression struct{ Name lexer.Token }

func (e *VariableExpression) Accept(visitor Visitor) { visitor.VisitVariableExpression(e) }
