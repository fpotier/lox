package runtime

import (
	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/lexer"
)

type LoxInstance struct {
	class  LoxClass
	fields map[string]ast.LoxValue
}

func NewLoxInstance(class LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]ast.LoxValue),
	}
}

func (i *LoxInstance) Get(name lexer.Token) ast.LoxValue {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	if method, ok := i.class.findMethod(name.Lexeme); ok {
		return method.Bind(i)
	}

	panic(NewUndefinedProperty(name.Line, name.Lexeme, i.class.name))
}

func (i *LoxInstance) Set(name lexer.Token, value ast.LoxValue) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) Kind() ast.Kind             { return ast.Instance }
func (i *LoxInstance) IsTruthy() bool             { return true }
func (i *LoxInstance) String() string             { return i.class.name + " instance" }
func (i *LoxInstance) Equals(_ ast.LoxValue) bool { return false }
