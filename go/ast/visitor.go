package ast

type Visitor interface {
	VisitBinaryExpression(*BinaryExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitUnaryExpression(*UnaryExpression)
	VisitVariableExpression(*VariableExpression)
	VisitAssignmentExpression(*AssignmentExpression)

	VisitExpressionStatement(*ExpressionStatement)
	VisitPrintStatement(*PrintStatement)
	VisitVariableStatement(*VariableStatement)
}
