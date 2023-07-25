package ast

import (
	"fmt"
)

type LoxValue interface {
	IsBoolean() bool
	IsNumber() bool
	IsString() bool
	IsTruthy() bool
	String() string
	Equals(v LoxValue) bool
}

type BooleanValue struct{ Value bool }

func NewBooleanValue(v bool) *BooleanValue { return &BooleanValue{Value: v} }
func (b BooleanValue) IsBoolean() bool     { return true }
func (b BooleanValue) IsNumber() bool      { return false }
func (b BooleanValue) IsString() bool      { return false }
func (b BooleanValue) IsTruthy() bool      { return b.Value }
func (b BooleanValue) String() string      { return fmt.Sprintf("%v", b.Value) }
func (b BooleanValue) Equals(v LoxValue) bool {
	if v.IsBoolean() {
		return b.Value == v.(*BooleanValue).Value
	}
	return false
}

type StringValue struct{ Value string }

func NewStringValue(v string) *StringValue { return &StringValue{Value: v} }
func (s StringValue) IsBoolean() bool      { return false }
func (s StringValue) IsNumber() bool       { return false }
func (s StringValue) IsString() bool       { return true }
func (s StringValue) IsTruthy() bool       { return true }
func (s StringValue) String() string       { return s.Value }
func (s StringValue) Equals(v LoxValue) bool {
	if v.IsString() {
		return s.Value == v.(*StringValue).Value
	}
	return false
}

type NumberValue struct{ Value float64 }

func NewNumberValue(v float64) *NumberValue { return &NumberValue{Value: v} }
func (n NumberValue) IsBoolean() bool       { return false }
func (n NumberValue) IsNumber() bool        { return true }
func (n NumberValue) IsString() bool        { return false }
func (n NumberValue) IsTruthy() bool        { return true }
func (n NumberValue) String() string        { return fmt.Sprintf("%v", n.Value) }
func (n NumberValue) Equals(v LoxValue) bool {
	if v.IsNumber() {
		return n.Value == v.(*NumberValue).Value
	}
	return false
}

type ObjectValue struct{ Value *map[string]LoxValue }

func (o ObjectValue) IsBoolean() bool { return false }
func (o ObjectValue) IsNumber() bool  { return false }
func (o ObjectValue) IsString() bool  { return false }
func (o ObjectValue) IsTruthy() bool  { return o.Value != nil }
func (o ObjectValue) String() string {
	// TODO: string representation of objects
	return "TODO"
}
func (o ObjectValue) Equals(_ LoxValue) bool {
	// TODO implement
	return false
}
