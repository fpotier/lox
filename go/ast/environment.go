package ast

import (
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Environment struct {
	values map[string]LoxValue
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]LoxValue)}
}

func (env *Environment) Define(name string, value LoxValue) {
	// Note: redeclaring a top-level (global) variable is allowed
	env.values[name] = value
}

func (env *Environment) Get(name lexer.Token) LoxValue {
	if val, ok := env.values[name.Lexeme]; ok {
		return val
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}

func (env *Environment) Assign(name lexer.Token, value LoxValue) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}
