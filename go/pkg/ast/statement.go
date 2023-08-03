package ast

import "github.com/fpotier/lox/go/pkg/lexer"

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		Statements: statements,
	}
}

func NewClassStatement(name lexer.Token, superclass *VariableExpression, methods []*FunctionStatement) *ClassStatement {
	return &ClassStatement{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{Expression: expression}
}

func NewFunctionStatement(name lexer.Token, parameters []lexer.Token, body []Statement) *FunctionStatement {
	return &FunctionStatement{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}
}

func NewIfStatment(condition Expression, thenCode Statement, elseCode Statement) *IfStatement {
	return &IfStatement{
		Condition: condition,
		ThenCode:  thenCode,
		ElseCode:  elseCode,
	}
}

func NewPrintStatement(expression Expression) *PrintStatement {
	return &PrintStatement{Expression: expression}
}

func NewReturnStatement(keyword lexer.Token, value Expression) *ReturnStatement {
	return &ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func NewVariableStatement(name lexer.Token, initializer Expression) *VariableStatement {
	return &VariableStatement{
		Name:        name,
		Initializer: initializer,
	}
}

func NewWhileStatement(condition Expression, body Statement) *WhileStatement {
	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) Accept(visitor Visitor) { visitor.VisitBlockStatement(s) }

type ClassStatement struct {
	Name       lexer.Token
	Superclass *VariableExpression
	Methods    []*FunctionStatement
}

func (s *ClassStatement) Accept(visitor Visitor) { visitor.VisitClassStatement(s) }

type ExpressionStatement struct {
	Expression Expression
}

func (s *ExpressionStatement) Accept(visitor Visitor) { visitor.VisitExpressionStatement(s) }

type FunctionStatement struct {
	Name       lexer.Token
	Parameters []lexer.Token
	Body       []Statement
}

func (s *FunctionStatement) Accept(visitor Visitor) { visitor.VisitFunctionStatement(s) }

type IfStatement struct {
	Condition Expression
	ThenCode  Statement
	ElseCode  Statement
}

func (s *IfStatement) Accept(visitor Visitor) { visitor.VisitIfStatement(s) }

type PrintStatement struct {
	Expression Expression
}

func (s *PrintStatement) Accept(visitor Visitor) { visitor.VisitPrintStatement(s) }

type ReturnStatement struct {
	Keyword lexer.Token
	Value   Expression
}

func (s *ReturnStatement) Accept(visitor Visitor) { visitor.VisitReturnStatement(s) }

type VariableStatement struct {
	Name        lexer.Token
	Initializer Expression
}

func (s *VariableStatement) Accept(visitor Visitor) { visitor.VisitVariableStatement(s) }

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (s *WhileStatement) Accept(visitor Visitor) { visitor.VisitWhileStatement(s) }
