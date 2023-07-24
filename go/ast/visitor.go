package ast

type Visitor interface {
	VisitBinaryExpression(*BinaryExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitUnaryExpression(*UnaryExpression)
	VisitVariableExpression(*VariableExpression)
	VisitAssignmentExpression(*AssignmentExpression)
	VisitLogicalExpression(*LogicalExpression)

	VisitExpressionStatement(*ExpressionStatement)
	VisitPrintStatement(*PrintStatement)
	VisitVariableStatement(*VariableStatement)
	VisitBlockStatement(*BlockStatement)
	VisitIfStatement(*IfStatement)
}
