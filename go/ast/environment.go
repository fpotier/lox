package ast

import (
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Environment struct {
	enclosing *Environment
	values    map[string]LoxValue
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]LoxValue),
	}
}

func NewSubEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]LoxValue),
	}
}

func (env *Environment) Define(name string, value LoxValue) {
	// Note: redeclaring a top-level (global) variable is allowed
	env.values[name] = value
}

func (env *Environment) Get(name lexer.Token) LoxValue {
	if val, ok := env.values[name.Lexeme]; ok {
		return val
	}

	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}

func (env *Environment) Assign(name lexer.Token, value LoxValue) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	if env.enclosing != nil {
		env.enclosing.Assign(name, value)
		return
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}
