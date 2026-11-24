package runtime

import (
	"github.com/fpotier/lox/go/pkg/lexer"
	"github.com/fpotier/lox/go/pkg/value"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]value.Value
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]value.Value),
	}
}

func (i *LoxInstance) Get(name lexer.Token) value.Value {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	if method, ok := i.class.findMethod(name.Lexeme); ok {
		return method.Bind(i)
	}

	panic(NewUndefinedProperty(name.Line, name.Lexeme, i.class.name))
}

func (i *LoxInstance) Set(name lexer.Token, value value.Value) {
	i.fields[name.Lexeme] = value
}

func (*LoxInstance) TypeName() string { return "instance" }
func (i *LoxInstance) String() string { return i.class.name + " instance" }
