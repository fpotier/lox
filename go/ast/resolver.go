package ast

import (
	"fmt"

	"github.com/fpotier/lox/go/lexer"
	"github.com/fpotier/lox/go/loxerror"
)

type FunctionType uint8

const (
	NoFunc FunctionType = iota
	Func
	Method
	Constructor
)

type ClassType uint8

const (
	NoClass ClassType = iota
	InClass
)

type Resolver struct {
	errorReporter    loxerror.ErrorReporter
	interpreter      *Interpreter
	scopes           []map[string]bool
	currentFnType    FunctionType
	currentClassType ClassType
}

func NewResolver(errorReporter loxerror.ErrorReporter, i *Interpreter) *Resolver {
	r := Resolver{
		errorReporter:    errorReporter,
		interpreter:      i,
		scopes:           make([]map[string]bool, 0),
		currentFnType:    NoFunc,
		currentClassType: NoClass,
	}

	return &r
}

func (r *Resolver) ResolveProgram(program []Statement) {
	for _, statement := range program {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) VisitGetExpression(e *GetExpression) {
	r.resolveExpression(e.Object)
}

func (r *Resolver) VisitSetExpression(e *SetExpression) {
	r.resolveExpression(e.Value)
	r.resolveExpression(e.Object)
}

func (r *Resolver) VisitVariableExpression(e *VariableExpression) {
	if len(r.scopes) > 0 {
		if declared, ok := r.scopes[len(r.scopes)-1][e.Name.Lexeme]; ok && !declared {
			r.errorReporter.Error(e.Name.Line, "Can't read local variable in its own initializer")
		}
	}

	r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitClassStatement(s *ClassStatement) {
	enclosingClass := r.currentClassType
	r.currentClassType = InClass

	r.declare(s.Name)
	r.define(s.Name)

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true
	for _, method := range s.Methods {
		fnType := Method
		if method.Name.Lexeme == "init" {
			fnType = Constructor
		}
		r.resolveFunction(*method, fnType)
	}
	r.endScope()

	r.currentClassType = enclosingClass
}

func (r *Resolver) VisitBlockStatement(s *BlockStatement) {
	r.beginScope()
	for _, statement := range s.Statements {
		r.resolveStatement(statement)
	}
	r.endScope()
}

func (r *Resolver) VisitVariableStatement(s *VariableStatement) {
	r.declare(s.Name)
	if s.Initializer != nil {
		r.resolveExpression(s.Initializer)
	}
	r.define(s.Name)
}

func (r *Resolver) VisitFunctionStatement(s *FunctionStatement) {
	r.declare(s.Name)
	r.define(s.Name)

	r.resolveFunction(*s, Func)
}

func (r *Resolver) VisitExpressionStatement(s *ExpressionStatement) {
	r.resolveExpression(s.Expression)
}

func (r *Resolver) VisitPrintStatement(s *PrintStatement) {
	r.resolveExpression(s.Expression)
}

func (r *Resolver) VisitIfStatement(s *IfStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.ThenCode)
	if s.ElseCode != nil {
		r.resolveStatement(s.ElseCode)
	}
}

func (r *Resolver) VisitWhileStatement(s *WhileStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.Body)
}

func (r *Resolver) VisitReturnStatement(s *ReturnStatement) {
	if r.currentFnType == NoFunc {
		r.errorReporter.Error(s.Keyword.Line, "Can't return from top-level code")
	}

	if s.Value != nil {
		if r.currentFnType == Constructor {
			r.errorReporter.Error(s.Keyword.Line, "Can't return from constructor")
		}

		r.resolveExpression(s.Value)
	}
}

func (r *Resolver) VisitBinaryExpression(e *BinaryExpression) {
	r.resolveExpression(e.LHS)
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitGroupingExpression(e *GroupingExpression) {
	r.resolveExpression(e.Expr)
}

func (r *Resolver) VisitLiteralExpression(*LiteralExpression) {}

func (r *Resolver) VisitUnaryExpression(e *UnaryExpression) {
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitLogicalExpression(e *LogicalExpression) {
	r.resolveExpression(e.LHS)
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitCallExpression(e *CallExpression) {
	r.resolveExpression(e.Callee)

	for _, argument := range e.Args {
		r.resolveExpression(argument)
	}
}

func (r *Resolver) VisitAssignmentExpression(e *AssignmentExpression) {
	r.resolveExpression(e.Value)
	r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitThisExpression(e *ThisExpression) {
	if r.currentClassType == NoClass {
		r.errorReporter.Error(e.Keyword.Line, "Can't use 'this' outside of a class")
	}

	r.resolveLocal(e, e.Keyword)
}

func (r *Resolver) resolveStatement(s Statement)   { s.Accept(r) }
func (r *Resolver) resolveExpression(e Expression) { e.Accept(r) }

func (r *Resolver) resolveLocal(e Expression, name lexer.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(e, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function FunctionStatement, kind FunctionType) {
	enclosingFunction := r.currentFnType
	r.currentFnType = kind

	r.beginScope()
	for _, parameter := range function.Parameters {
		r.declare(parameter)
		r.define(parameter)
	}
	for _, statement := range function.Body {
		r.resolveStatement(statement)
	}
	r.endScope()

	r.currentFnType = enclosingFunction
}

func (r *Resolver) beginScope() { r.scopes = append(r.scopes, make(map[string]bool)) }
func (r *Resolver) endScope()   { r.scopes = r.scopes[:len(r.scopes)-1] }

func (r *Resolver) declare(name lexer.Token) {
	if len(r.scopes) > 0 {
		if _, ok := r.scopes[len(r.scopes)-1][name.Lexeme]; ok {
			r.errorReporter.Error(name.Line, fmt.Sprintf("Variable '%s' already declared in this scope", name.Lexeme))
		}
		r.scopes[len(r.scopes)-1][name.Lexeme] = false
	}
}

func (r *Resolver) define(name lexer.Token) {
	if len(r.scopes) > 0 {
		r.scopes[len(r.scopes)-1][name.Lexeme] = true
	}
}
