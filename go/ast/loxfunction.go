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

func (f NativeFunction) IsBoolean() bool        { return false }
func (f NativeFunction) IsNumber() bool         { return false }
func (f NativeFunction) IsString() bool         { return false }
func (f NativeFunction) IsTruthy() bool         { return false }
func (f NativeFunction) String() string         { return fmt.Sprintf("<native function> %s", f.name) }
func (f NativeFunction) Equals(_ LoxValue) bool { return false }
func (f NativeFunction) Call(i *Interpreter, arguments []LoxValue) LoxValue {
	return f.code(i, arguments)
}
func (f NativeFunction) Arity() int { return f.arity }

type LoxFunction struct {
	Declaration *FunctionStatement
	Closure     *Environment
}

func NewLoxFunction(declaration *FunctionStatement, closure *Environment) *LoxFunction {
	return &LoxFunction{Declaration: declaration, Closure: closure}
}
func (f LoxFunction) IsBoolean() bool        { return false }
func (f LoxFunction) IsNumber() bool         { return false }
func (f LoxFunction) IsString() bool         { return false }
func (f LoxFunction) IsTruthy() bool         { return false }
func (f LoxFunction) String() string         { return fmt.Sprintf("<function> %s", f.Declaration.Name.Lexeme) }
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
				returnValue = rt
				return
			}
			panic(r)
		}
	}()

	i.executeBlock(f.Declaration.Body, environment)

	return returnValue
}
func (f LoxFunction) Arity() int { return len(f.Declaration.Parameters) }
