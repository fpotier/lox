package ast

import (
	"github.com/fpotier/lox/go/lexer"
	"github.com/fpotier/lox/go/loxerror"
)

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:    name,
		methods: methods,
	}
}

func (c *LoxClass) Kind() Kind             { return Class }
func (c *LoxClass) IsTruthy() bool         { return true }
func (c *LoxClass) String() string         { return c.name }
func (c *LoxClass) Equals(v LoxValue) bool { return c == v }

func (c *LoxClass) Call(i *Interpreter, arguments []LoxValue) LoxValue {
	instance := NewLoxInstance(*c)
	if constructor, ok := c.findMethod("init"); ok {
		constructor.Bind(instance).Call(i, arguments)
	}

	return instance
}

func (c *LoxClass) Arity() int {
	if constructor, ok := c.findMethod("init"); ok {
		return constructor.Arity()
	}
	return 0
}

func (c *LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if method, ok := c.methods[name]; ok {
		return method, true
	}
	return nil, false
}

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

	if method, ok := i.class.findMethod(name.Lexeme); ok {
		return method.Bind(i)
	}

	panic(loxerror.RuntimeError{Message: "Undefined property " + name.Lexeme})
}

func (i *LoxInstance) Set(name lexer.Token, value LoxValue) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) Kind() Kind             { return Instance }
func (i *LoxInstance) IsTruthy() bool         { return true }
func (i *LoxInstance) String() string         { return i.class.name + " instance" }
func (i *LoxInstance) Equals(_ LoxValue) bool { return false }
