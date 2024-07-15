package ast

type Expression interface {
	Accept(visitor Visitor)
}

type Statement interface {
	Accept(visitor Visitor)
}

type Visitor interface {
	VisitAssignmentExpression(assignementExpression *AssignmentExpression)
	VisitBinaryExpression(binaryExpression *BinaryExpression)
	VisitCallExpression(callExpression *CallExpression)
	VisitGetExpression(getExpression *GetExpression)
	VisitGroupingExpression(groupingExpression *GroupingExpression)
	VisitLiteralExpression(literalExpression *LiteralExpression)
	VisitLogicalExpression(logicalExpression *LogicalExpression)
	VisitSetExpression(setExpression *SetExpression)
	VisitSuperExpression(superExpression *SuperExpression)
	VisitThisExpression(thisExpression *ThisExpression)
	VisitUnaryExpression(unaryExpression *UnaryExpression)
	VisitVariableExpression(variableExpression *VariableExpression)

	VisitBlockStatement(blockStatement *BlockStatement)
	VisitClassStatement(classStatement *ClassStatement)
	VisitExpressionStatement(expressionStatement *ExpressionStatement)
	VisitFunctionStatement(functionStatement *FunctionStatement)
	VisitIfStatement(ifStatement *IfStatement)
	VisitPrintStatement(printStatement *PrintStatement)
	VisitReturnStatement(returnStatement *ReturnStatement)
	VisitVariableStatement(variableStatement *VariableStatement)
	VisitWhileStatement(whileStatement *WhileStatement)
}
