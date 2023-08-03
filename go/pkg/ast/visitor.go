package ast

type Expression interface {
	Accept(visitor Visitor)
}

type Statement interface {
	Accept(visitor Visitor)
}

type Visitor interface {
	VisitAssignmentExpression(*AssignmentExpression)
	VisitBinaryExpression(*BinaryExpression)
	VisitCallExpression(*CallExpression)
	VisitGetExpression(*GetExpression)
	VisitGroupingExpression(*GroupingExpression)
	VisitLiteralExpression(*LiteralExpression)
	VisitLogicalExpression(*LogicalExpression)
	VisitSetExpression(*SetExpression)
	VisitSuperExpression(*SuperExpression)
	VisitThisExpression(*ThisExpression)
	VisitUnaryExpression(*UnaryExpression)
	VisitVariableExpression(*VariableExpression)

	VisitBlockStatement(*BlockStatement)
	VisitClassStatement(*ClassStatement)
	VisitExpressionStatement(*ExpressionStatement)
	VisitFunctionStatement(*FunctionStatement)
	VisitIfStatement(*IfStatement)
	VisitPrintStatement(*PrintStatement)
	VisitReturnStatement(*ReturnStatement)
	VisitVariableStatement(*VariableStatement)
	VisitWhileStatement(*WhileStatement)
}
