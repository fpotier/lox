package runtime

import (
	"fmt"

	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/value"
)

type LoxFunction struct {
	Declaration   *ast.FunctionStatement
	Closure       *Environment
	isConstructor bool
	className     string
}

func NewLoxFunction(declaration *ast.FunctionStatement, closure *Environment, isConstructor bool) *LoxFunction {
	return &LoxFunction{Declaration: declaration, Closure: closure, isConstructor: isConstructor, className: ""}
}
func (f *LoxFunction) setClassName(className string) { f.className = className }
func (*LoxFunction) TypeName() string                { return "function" }
func (f *LoxFunction) String() string                { return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme) }
func (f *LoxFunction) Name() string {
	if len(f.className) == 0 {
		return f.Declaration.Name.Lexeme
	}

	return f.className + "::" + f.Declaration.Name.Lexeme
}

func (f *LoxFunction) Call(i *Interpreter, arguments []value.Value) (returnValue value.Value) {
	environment := NewSubEnvironment(f.Closure)
	for i := range f.Declaration.Parameters {
		environment.Define(f.Declaration.Parameters[i].Lexeme, arguments[i])
	}

	hasReturned := i.hasReturned
	i.hasReturned = false
	i.executeBlock(f.Declaration.Body, environment)
	i.hasReturned = hasReturned

	if f.isConstructor {
		return f.Closure.GetAt(0, "this")
	}

	return i.Value
}
func (f *LoxFunction) Arity() int { return len(f.Declaration.Parameters) }

func (f *LoxFunction) Bind(this *LoxInstance) *LoxFunction {
	environment := NewSubEnvironment(f.Closure)
	environment.Define("this", this)

	boundFunction := NewLoxFunction(f.Declaration, environment, f.isConstructor)

	return boundFunction
}
