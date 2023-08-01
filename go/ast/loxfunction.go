package ast

import "fmt"

const maxArgs = 255

type LoxCallable interface {
	Call(*Interpreter, []LoxValue) LoxValue
	Arity() int
}

type NativeFunction struct {
	name  string
	arity int
	code  func(*Interpreter, []LoxValue) LoxValue
}

func (f NativeFunction) Kind() Kind             { return NativeFunc }
func (f NativeFunction) IsTruthy() bool         { return true }
func (f NativeFunction) String() string         { return "<native fn>" }
func (f NativeFunction) Equals(_ LoxValue) bool { return false }
func (f NativeFunction) Call(i *Interpreter, arguments []LoxValue) LoxValue {
	return f.code(i, arguments)
}
func (f NativeFunction) Arity() int { return f.arity }

type LoxFunction struct {
	Declaration   *FunctionStatement
	Closure       *Environment
	isConstructor bool
}

func NewLoxFunction(declaration *FunctionStatement, closure *Environment, isConstructor bool) *LoxFunction {
	return &LoxFunction{Declaration: declaration, Closure: closure, isConstructor: isConstructor}
}
func (f LoxFunction) Kind() Kind             { return Function }
func (f LoxFunction) IsTruthy() bool         { return true }
func (f LoxFunction) String() string         { return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme) }
func (f LoxFunction) Equals(_ LoxValue) bool { return false }
func (f LoxFunction) Call(i *Interpreter, arguments []LoxValue) (returnValue LoxValue) {
	environment := NewSubEnvironment(f.Closure)
	for i := range f.Declaration.Parameters {
		environment.Define(f.Declaration.Parameters[i].Lexeme, arguments[i])
	}

	// If no return statement is executed 'nil' is returned
	returnValue = NewNilValue()

	defer func() {
		if r := recover(); r != nil {
			if rt, ok := r.(LoxValue); ok {
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
