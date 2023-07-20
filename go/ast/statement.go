package ast

import (
	"github.com/fpotier/crafting-interpreters/go/lexer"
)

type Statement interface {
	Accept(visitor Visitor)
}

type ExpressionStatement struct {
	Expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{Expression: expression}
}

func (s *ExpressionStatement) Accept(visitor Visitor) {
	visitor.VisitExpressionStatement(s)
}

type PrintStatement struct {
	Expression Expression
}

func NewPrintStatement(expression Expression) *PrintStatement {
	return &PrintStatement{Expression: expression}
}

func (s *PrintStatement) Accept(visitor Visitor) {
	visitor.VisitPrintStatement(s)
}

type VariableStatement struct {
	Name        lexer.Token
	Initializer Expression
}

func NewVariableStatement(name lexer.Token, initializer Expression) *VariableStatement {
	return &VariableStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func (s *VariableStatement) Accept(visitor Visitor) {
	visitor.VisitVariableStatement(s)
}

type BlockStatement struct {
	Statements []Statement
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		Statements: statements,
	}
}

func (s *BlockStatement) Accept(visitor Visitor) {
	visitor.VisitBlockStatement(s)
}
