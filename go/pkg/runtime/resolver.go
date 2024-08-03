package runtime

import (
	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/lexer"
	"github.com/fpotier/lox/go/pkg/loxerror"
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
	InSubClass
)

type Resolver struct {
	errorFormatter   loxerror.ErrorFormatter
	interpreter      *Interpreter
	scopes           []map[string]bool
	currentFnType    FunctionType
	currentClassType ClassType
}

func NewResolver(errorFormatter loxerror.ErrorFormatter, i *Interpreter) *Resolver {
	r := Resolver{
		errorFormatter:   errorFormatter,
		interpreter:      i,
		scopes:           make([]map[string]bool, 0),
		currentFnType:    NoFunc,
		currentClassType: NoClass,
	}

	return &r
}

func (r *Resolver) ResolveProgram(program []ast.Statement) {
	for _, statement := range program {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) VisitAssignmentExpression(e *ast.AssignmentExpression) {
	r.resolveExpression(e.Value)
	r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitBinaryExpression(e *ast.BinaryExpression) {
	r.resolveExpression(e.LHS)
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitCallExpression(e *ast.CallExpression) {
	r.resolveExpression(e.Callee)

	for _, argument := range e.Args {
		r.resolveExpression(argument)
	}
}

func (r *Resolver) VisitGetExpression(e *ast.GetExpression) {
	r.resolveExpression(e.Object)
}

func (r *Resolver) VisitGroupingExpression(e *ast.GroupingExpression) {
	r.resolveExpression(e.Expr)
}

func (r *Resolver) VisitLiteralExpression(*ast.LiteralExpression) {}

func (r *Resolver) VisitLogicalExpression(e *ast.LogicalExpression) {
	r.resolveExpression(e.LHS)
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitSetExpression(e *ast.SetExpression) {
	r.resolveExpression(e.Value)
	r.resolveExpression(e.Object)
}

func (r *Resolver) VisitSuperExpression(e *ast.SuperExpression) {
	if r.currentClassType == NoClass {
		r.errorFormatter.PushError(NewInvalidSuper(e.Keyword.Line, "outside of a class"))
	} else if r.currentClassType != InSubClass {
		r.errorFormatter.PushError(NewInvalidSuper(e.Keyword.Line, "in a class with no superclass"))
	}

	r.resolveLocal(e, e.Keyword)
}

func (r *Resolver) VisitThisExpression(e *ast.ThisExpression) {
	if r.currentClassType == NoClass {
		r.errorFormatter.PushError(NewInvalidThis(e.Keyword.Line, "outside of a class"))
	}

	r.resolveLocal(e, e.Keyword)
}

func (r *Resolver) VisitUnaryExpression(e *ast.UnaryExpression) {
	r.resolveExpression(e.RHS)
}

func (r *Resolver) VisitVariableExpression(e *ast.VariableExpression) {
	if len(r.scopes) > 0 {
		if declared, ok := r.scopes[len(r.scopes)-1][e.Name.Lexeme]; ok && !declared {
			r.errorFormatter.PushError(NewUninitializedRead(e.Name.Line))
		}
	}

	r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitBlockStatement(s *ast.BlockStatement) {
	r.beginScope()
	for _, statement := range s.Statements {
		r.resolveStatement(statement)
	}
	r.endScope()
}

func (r *Resolver) VisitClassStatement(s *ast.ClassStatement) {
	enclosingClass := r.currentClassType
	r.currentClassType = InClass

	r.declare(s.Name)
	r.define(s.Name)

	if s.Superclass != nil {
		r.currentClassType = InSubClass
		if s.Name.Lexeme == s.Superclass.Name.Lexeme {
			r.errorFormatter.PushError(NewInvalidInheritance(s.Superclass.Name.Line, "A class can't inherit from itself"))
		}
		r.resolveExpression(s.Superclass)

		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

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

	if s.Superclass != nil {
		r.endScope()
	}

	r.currentClassType = enclosingClass
}

func (r *Resolver) VisitExpressionStatement(s *ast.ExpressionStatement) {
	r.resolveExpression(s.Expression)
}

func (r *Resolver) VisitFunctionStatement(s *ast.FunctionStatement) {
	r.declare(s.Name)
	r.define(s.Name)

	r.resolveFunction(*s, Func)
}

func (r *Resolver) VisitIfStatement(s *ast.IfStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.ThenCode)
	if s.ElseCode != nil {
		r.resolveStatement(s.ElseCode)
	}
}

func (r *Resolver) VisitPrintStatement(s *ast.PrintStatement) {
	r.resolveExpression(s.Expression)
}

func (r *Resolver) VisitReturnStatement(s *ast.ReturnStatement) {
	if r.currentFnType == NoFunc {
		r.errorFormatter.PushError(NewInvalidReturn(s.Keyword.Line, "top-level code"))
	}

	if s.Value != nil {
		if r.currentFnType == Constructor {
			r.errorFormatter.PushError(NewInvalidReturn(s.Keyword.Line, "constructor"))
		}

		r.resolveExpression(s.Value)
	}
}

func (r *Resolver) VisitVariableStatement(s *ast.VariableStatement) {
	r.declare(s.Name)
	if s.Initializer != nil {
		r.resolveExpression(s.Initializer)
	}
	r.define(s.Name)
}

func (r *Resolver) VisitWhileStatement(s *ast.WhileStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.Body)
}

func (r *Resolver) resolveStatement(s ast.Statement)   { s.Accept(r) }
func (r *Resolver) resolveExpression(e ast.Expression) { e.Accept(r) }

func (r *Resolver) resolveLocal(e ast.Expression, name lexer.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(e, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function ast.FunctionStatement, kind FunctionType) {
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
			r.errorFormatter.PushError(NewVariableRedeclaration(name.Line, name.Lexeme))
		}
		r.scopes[len(r.scopes)-1][name.Lexeme] = false
	}
}

func (r *Resolver) define(name lexer.Token) {
	if len(r.scopes) > 0 {
		r.scopes[len(r.scopes)-1][name.Lexeme] = true
	}
}
