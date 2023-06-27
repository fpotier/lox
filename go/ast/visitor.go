package ast

type Visitor interface {
	VisitBinaryExpression(binaryExpression *BinaryExpression)
	VisitGroupingExpression(groupingExpression *GroupingExpression)
	VisitLiteralExpression(literalExpression *LiteralExpression)
	VisitUnaryExpression(unaryExpression *UnaryExpression)

	VisitExpressionStatement(expressionStatement *ExpressionStatement)
	VisitPrintStatement(printStatement *PrintStatement)
}
