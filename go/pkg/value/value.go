// Package value defines the runtime representation of Lox values.
//
// Primitive values are plain Go named types, so they are compared with ==
// and inspected with type assertions or type switches. Reference values
// (functions, classes, instances) live in the runtime package and are
// always pointers, so == on them is identity.
package value

import (
	"fmt"
	"strconv"
)

// Value is any Lox value.
type Value interface {
	// String returns the representation printed by the "print" statement.
	String() string
	// TypeName returns the user-facing name of the value's type, used in error messages.
	TypeName() string
}

type (
	Bool   bool
	Number float64
	String string
	Nil    struct{}
)

func (b Bool) String() string   { return strconv.FormatBool(bool(b)) }
func (Bool) TypeName() string   { return "boolean" }
func (n Number) String() string { return fmt.Sprint(float64(n)) }
func (Number) TypeName() string { return "number" }
func (s String) String() string { return string(s) }
func (String) TypeName() string { return "string" }
func (Nil) String() string      { return "nil" }
func (Nil) TypeName() string    { return "nil" }

// IsTruthy reports whether v is considered true in a boolean context:
// nil and false are falsy, everything else is truthy.
func IsTruthy(v Value) bool {
	switch v := v.(type) {
	case Nil:
		return false
	case Bool:
		return bool(v)
	default:
		return true
	}
}
