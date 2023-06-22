package ast

import "fmt"

type LoxValue interface {
	IsTruthy() bool
	String() string
}

type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) IsTruthy() bool {
	return b.Value
}

func (b *BooleanValue) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

type StringValue struct {
	Value string
}

func (s *StringValue) IsTruthy() bool {
	return true
}

func (s *StringValue) String() string {
	return s.Value
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) IsTruthy() bool {
	return true
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type ObjectValue struct {
	Value *map[string]LoxValue
}

func (o *ObjectValue) IsTruthy() bool {
	return o.Value != nil
}

func (o *ObjectValue) String() string {
	// TODO: string representation of objects
	return "TODO"
}
