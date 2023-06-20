package ast

import (
	"fmt"
	"strings"
)

type PrinterVisitor struct {
}

func (visitor *PrinterVisitor) Print(expr Expression) string {
	return expr.Accept(visitor).(string)
}

func (visitor *PrinterVisitor) VisitBinaryExpression(binaryExpression *BinaryExpression) interface{} {
	return parenthesize(binaryExpression.operator.Lexeme, binaryExpression.lhs, binaryExpression.rhs)
}

func (visitor *PrinterVisitor) VisitGroupingExpression(groupingExpression *GroupingExpression) interface{} {
	return parenthesize("group", groupingExpression.expr)
}

func (visitor *PrinterVisitor) VisitLiteralExpression(literalExpression *LiteralExpression) interface{} {
	if literalExpression.value.IsNumber {
		return fmt.Sprintf("%v", literalExpression.value.NumberValue)
	} else if literalExpression.value.IsString {
		return fmt.Sprintf("%v", literalExpression.value.StringValue)
	} else {
		return "nil"
	}
}

func (visitor *PrinterVisitor) VisitUnaryExpression(unaryExpression *UnaryExpression) interface{} {
	return parenthesize(unaryExpression.operator.Lexeme, unaryExpression.rhs)
}

func parenthesize(name string, expressions ...Expression) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range expressions {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(&PrinterVisitor{}).(string))
	}
	builder.WriteString(")")

	return builder.String()
}
