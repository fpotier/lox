package runtime

import (
	"fmt"

	"github.com/fpotier/lox/go/pkg/ast"
)

type LoxFunction struct {
	Declaration   *ast.FunctionStatement
	Closure       *Environment
	isConstructor bool
}

func NewLoxFunction(declaration *ast.FunctionStatement, closure *Environment, isConstructor bool) *LoxFunction {
	return &LoxFunction{Declaration: declaration, Closure: closure, isConstructor: isConstructor}
}
func (f LoxFunction) Kind() ast.Kind             { return ast.Function }
func (f LoxFunction) IsTruthy() bool             { return true }
func (f LoxFunction) String() string             { return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme) }
func (f LoxFunction) Equals(_ ast.LoxValue) bool { return false }
func (f LoxFunction) Call(i *Interpreter, arguments []ast.LoxValue) (returnValue ast.LoxValue) {
	environment := NewSubEnvironment(f.Closure)
	for i := range f.Declaration.Parameters {
		environment.Define(f.Declaration.Parameters[i].Lexeme, arguments[i])
	}

	// If no return statement is executed 'nil' is returned
	returnValue = ast.NewNilValue()

	defer func() {
		if r := recover(); r != nil {
			if rt, ok := r.(ast.LoxValue); ok {
				if f.isConstructor {
					returnValue = f.Closure.GetAt(0, "this")
				} else {
					returnValue = rt
				}
				return
			}
			panic(r)
		}
	}()

	i.executeBlock(f.Declaration.Body, environment)

	if f.isConstructor {
		return f.Closure.GetAt(0, "this")
	}

	return returnValue
}
func (f LoxFunction) Arity() int { return len(f.Declaration.Parameters) }

func (f *LoxFunction) Bind(this *LoxInstance) *LoxFunction {
	environment := NewSubEnvironment(f.Closure)
	environment.Define("this", this)

	boundFunction := NewLoxFunction(f.Declaration, environment, f.isConstructor)

	if f.isConstructor {
		fmt.Println(f.Arity() == boundFunction.Arity())
	}

	return boundFunction
}
