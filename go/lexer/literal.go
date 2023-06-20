package lexer

import "fmt"

type Literal struct {
	IsString    bool
	StringValue string
	IsNumber    bool
	NumberValue float64
}

func (literal *Literal) String() string {
	if literal.IsNumber {
		return fmt.Sprintf("(value:%v)", literal.NumberValue)
	} else if literal.IsString {
		return fmt.Sprintf("(value:'%v')", literal.StringValue)
	} else {
		return ""
	}
}
