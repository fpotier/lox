package ast

type Visitor interface {
	VisitBinaryExpression(*BinaryExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitUnaryExpression(*UnaryExpression)
	VisitVariableExpression(*VariableExpression)
	VisitAssignmentExpression(*AssignmentExpression)
	VisitLogicalExpression(*LogicalExpression)
	VisitCallExpression(*CallExpression)

	VisitExpressionStatement(*ExpressionStatement)
	VisitPrintStatement(*PrintStatement)
	VisitVariableStatement(*VariableStatement)
	VisitBlockStatement(*BlockStatement)
	VisitIfStatement(*IfStatement)
	VisitWhileStatement(*WhileStatement)
	VisitFunctionStatement(*FunctionStatement)
	VisitReturnStatement(*ReturnStatement)
}
