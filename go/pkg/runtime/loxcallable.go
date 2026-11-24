package runtime

import (
	"github.com/fpotier/lox/go/pkg/value"
)

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []value.Value) value.Value
	Arity() int
	Name() string
}

type CallableCode func(*Interpreter, []value.Value) value.Value
