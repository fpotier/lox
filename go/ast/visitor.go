package ast

type Visitor interface {
	VisitBinaryExpression(binaryExpression *BinaryExpression) interface{}
	VisitGroupingExpression(groupingExpression *GroupingExpression) interface{}
	VisitLiteralExpression(literalExpression *LiteralExpression) interface{}
	VisitUnaryExpression(unaryExpression *UnaryExpression) interface{}
}
