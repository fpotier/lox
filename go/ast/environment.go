package ast

import (
	"fmt"

	"github.com/fpotier/lox/go/lexer"
	"github.com/fpotier/lox/go/loxerror"
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

func (e *Environment) GetAt(distance int, name string) LoxValue {
	return e.ancestor(distance).symbols[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
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

func (e *Environment) AssignAt(distance int, name lexer.Token, value LoxValue) {
	e.ancestor(distance).symbols[name.Lexeme] = value
}
