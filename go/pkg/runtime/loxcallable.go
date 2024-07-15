package runtime

import "github.com/fpotier/lox/go/pkg/ast"

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []ast.LoxValue) ast.LoxValue
	Arity() int
}

type CallableCode func(*Interpreter, []ast.LoxValue) ast.LoxValue
