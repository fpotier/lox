package ast

import (
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Environment struct {
	enclosing *Environment
	symbols   map[string]LoxValue
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		symbols:   make(map[string]LoxValue),
	}
}

func NewSubEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		symbols:   make(map[string]LoxValue),
	}
}

func (e *Environment) Define(name string, value LoxValue) {
	// Note: redeclaring a top-level (global) variable is allowed
	e.symbols[name] = value
}

func (e *Environment) Get(name lexer.Token) LoxValue {
	if val, ok := e.symbols[name.Lexeme]; ok {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}

func (e *Environment) Assign(name lexer.Token, value LoxValue) {
	if _, ok := e.symbols[name.Lexeme]; ok {
		e.symbols[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	panic(loxerror.RuntimeError{Message: fmt.Sprintf("Undefined variable '%v'", name.Lexeme)})
}
