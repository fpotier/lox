package ast

import (
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass { return &LoxClass{name: name} }

func (c *LoxClass) IsBoolean() bool        { return false }
func (c *LoxClass) IsNumber() bool         { return false }
func (c *LoxClass) IsString() bool         { return false }
func (c *LoxClass) IsTruthy() bool         { return false }
func (c *LoxClass) String() string         { return c.name }
func (c *LoxClass) Equals(_ LoxValue) bool { return false }

func (c *LoxClass) Call(i *Interpreter, arguments []LoxValue) LoxValue {
	return NewLoxInstance(*c)
}

func (c *LoxClass) Arity() int { return 0 }

type LoxInstance struct {
	class  LoxClass
	fields map[string]LoxValue
}

func NewLoxInstance(class LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]LoxValue),
	}
}

func (i *LoxInstance) Get(name lexer.Token) LoxValue {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	panic(loxerror.RuntimeError{Message: "Undefined property " + name.Lexeme})
}

func (i *LoxInstance) Set(name lexer.Token, value LoxValue) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) IsBoolean() bool        { return false }
func (i *LoxInstance) IsNumber() bool         { return false }
func (i *LoxInstance) IsString() bool         { return false }
func (i *LoxInstance) IsTruthy() bool         { return false }
func (i *LoxInstance) String() string         { return i.class.name + " instance" }
func (i *LoxInstance) Equals(_ LoxValue) bool { return false }
