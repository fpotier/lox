package visitor

import (
	"strings"

	"github.com/fpotier/crafting-interpreters/go/ast"
)

type RPNPrinter struct {
	builder strings.Builder
}

func (visitor *RPNPrinter) String(expr ast.Expression) string {
	expr.Accept(visitor)
	return visitor.builder.String()
}

func (visitor *RPNPrinter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	binaryExpression.Lhs.Accept(visitor)
	visitor.builder.WriteString(" ")
	binaryExpression.Rhs.Accept(visitor)
	visitor.builder.WriteString(" ")
	visitor.builder.WriteString(binaryExpression.Operator.Lexeme)
}

func (visitor *RPNPrinter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	groupingExpression.Expr.Accept(visitor)
}

func (visitor *RPNPrinter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	visitor.builder.WriteString(literalExpression.Value.String())
}

func (visitor *RPNPrinter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	visitor.builder.WriteString(unaryExpression.Operator.Lexeme)
	unaryExpression.Rhs.Accept(visitor)
}
