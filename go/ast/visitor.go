package ast

type Visitor interface {
	VisitBinaryExpression(binaryExpression *BinaryExpression)
	VisitGroupingExpression(groupingExpression *GroupingExpression)
	VisitLiteralExpression(literalExpression *LiteralExpression)
	VisitUnaryExpression(unaryExpression *UnaryExpression)
}
