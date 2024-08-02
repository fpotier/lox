package runtime

import (
	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/lexer"
)

type Environment struct {
	enclosing *Environment
	symbols   map[string]ast.LoxValue
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		symbols:   make(map[string]ast.LoxValue),
	}
}

func NewSubEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		symbols:   make(map[string]ast.LoxValue),
	}
}

func (e *Environment) Define(name string, value ast.LoxValue) {
	// Note: redeclaring a top-level (global) variable is allowed
	e.symbols[name] = value
}

func (e *Environment) Get(name lexer.Token) ast.LoxValue {
	if val, ok := e.symbols[name.Lexeme]; ok {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	panic(NewUndefinedVariable(name.Line, name.Lexeme))
}

func (e *Environment) GetAt(distance int, name string) ast.LoxValue {
	return e.ancestor(distance).symbols[name]
}

func (e *Environment) Assign(name lexer.Token, value ast.LoxValue) {
	if _, ok := e.symbols[name.Lexeme]; ok {
		e.symbols[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	panic(NewUndefinedVariable(name.Line, name.Lexeme))
}

func (e *Environment) AssignAt(distance int, name lexer.Token, value ast.LoxValue) {
	e.ancestor(distance).symbols[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}
