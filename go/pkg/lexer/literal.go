package lexer

import "fmt"

type Literal interface {
	String() string
}

type StringLiteral struct {
	Value string
}

func (l *StringLiteral) String() string {
	return l.Value
}

type NumberLiteral struct {
	Value float64
}

func (l *NumberLiteral) String() string {
	return fmt.Sprintf("%v", l.Value)
}
