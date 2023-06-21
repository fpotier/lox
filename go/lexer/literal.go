package lexer

import "fmt"

type Literal interface {
	String() string
}

type StringLiteral struct {
	Value string
}

func (literal *StringLiteral) String() string {
	return literal.Value
}

type NumberLiteral struct {
	Value float64
}

func (literal *NumberLiteral) String() string {
	return fmt.Sprintf("%v", literal.Value)
}
