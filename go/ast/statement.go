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

type IfStatement struct {
	Condition Expression
	ThenCode  Statement
	ElseCode  Statement
}

func NewIfStatment(condition Expression, thenCode Statement, elseCode Statement) *IfStatement {
	return &IfStatement{
		Condition: condition,
		ThenCode:  thenCode,
		ElseCode:  elseCode,
	}
}

func (s *IfStatement) Accept(visitor Visitor) {
	visitor.VisitIfStatement(s)
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func NewWhileStatement(condition Expression, body Statement) *WhileStatement {
	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (s *WhileStatement) Accept(visitor Visitor) {
	visitor.VisitWhileStatement(s)
}

type FunctionStatement struct {
	Name       lexer.Token
	Parameters []lexer.Token
	Body       []Statement
}

func NewFunctionStatement(name lexer.Token, parameters []lexer.Token, body []Statement) *FunctionStatement {
	return &FunctionStatement{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}
}

func (s *FunctionStatement) Accept(visitor Visitor) {
	visitor.VisitFunctionStatement(s)
}
