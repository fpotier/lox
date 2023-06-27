package visitor

import (
	"strings"
)

type LispPrinter struct {
	builder strings.Builder
}

/*
func (visitor *LispPrinter) String(expr ast.Expression) string {
	expr.Accept(visitor)
	return visitor.builder.String()
}

func (visitor *LispPrinter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	visitor.parenthesize(binaryExpression.Operator.Lexeme, binaryExpression.Lhs, binaryExpression.Rhs)
}

func (visitor *LispPrinter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	visitor.parenthesize("group", groupingExpression.Expr)
}

func (visitor *LispPrinter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	visitor.builder.WriteString(literalExpression.LoxValue().String())
}

func (visitor *LispPrinter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	visitor.parenthesize(unaryExpression.Operator.Lexeme, unaryExpression.Rhs)
}

func (visitor *LispPrinter) parenthesize(name string, expressions ...ast.Expression) {
	visitor.builder.WriteString("(")
	visitor.builder.WriteString(name)
	for _, expr := range expressions {
		visitor.builder.WriteString(" ")
		expr.Accept(visitor)
	}
	visitor.builder.WriteString(")")
}
*/
