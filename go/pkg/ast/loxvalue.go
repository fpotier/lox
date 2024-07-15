package ast

import (
	"fmt"
	"strconv"
)

type Kind uint8

const (
	Boolean Kind = iota
	String
	Number
	Nil
	NativeFunc
	Function
	Class
	Instance
)

type LoxValue interface {
	Kind() Kind
	IsTruthy() bool
	String() string
	Equals(v LoxValue) bool
}

type BooleanValue struct{ Value bool }

var True = &BooleanValue{Value: true}
var False = &BooleanValue{Value: false}

func NewBooleanValue(v bool) *BooleanValue {
	if v {
		return True
	}
	return False
}
func (b BooleanValue) Kind() Kind     { return Boolean }
func (b BooleanValue) IsTruthy() bool { return b.Value }
func (b BooleanValue) String() string { return strconv.FormatBool(b.Value) }
func (b BooleanValue) Equals(v LoxValue) bool {
	if v, ok := v.(*BooleanValue); ok {
		return b.Value == v.Value
	}
	return false
}

type StringValue struct{ Value string }

func NewStringValue(v string) *StringValue { return &StringValue{Value: v} }
func (s StringValue) Kind() Kind           { return String }
func (s StringValue) IsTruthy() bool       { return true }
func (s StringValue) String() string       { return s.Value }
func (s StringValue) Equals(v LoxValue) bool {
	if v, ok := v.(*StringValue); ok {
		return s.Value == v.Value
	}
	return false
}

type NumberValue struct{ Value float64 }

func NewNumberValue(v float64) *NumberValue { return &NumberValue{Value: v} }
func (n NumberValue) Kind() Kind            { return Number }
func (n NumberValue) IsTruthy() bool        { return true }
func (n NumberValue) String() string        { return fmt.Sprintf("%v", n.Value) }
func (n NumberValue) Equals(v LoxValue) bool {
	if v, ok := v.(*NumberValue); ok {
		return n.Value == v.Value
	}
	return false
}

type NilValue struct{}

var NilVal = &NilValue{}

func NewNilValue() *NilValue      { return NilVal }
func (n NilValue) Kind() Kind     { return Nil }
func (n NilValue) IsTruthy() bool { return false }
func (n NilValue) String() string { return "nil" }
func (n NilValue) Equals(v LoxValue) bool {
	if _, ok := v.(*NilValue); ok {
		return true
	}
	return false
}
