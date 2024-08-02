package ast

import (
	"fmt"
	"io"
	"strings"
)

type AstPrinter struct {
	output          io.Writer
	identationLevel int
	tabSize         int
}

func NewAstPrinter(outputStream io.Writer, tabSize uint) AstPrinter {
	return AstPrinter{
		output:          outputStream,
		identationLevel: 0,
		tabSize:         2,
	}
}

func (astPrinter *AstPrinter) Dump(statements []Statement) {
	for _, statement := range statements {
		statement.Accept(astPrinter)
	}
}

func (astPrinter AstPrinter) VisitAssignmentExpression(assignementExpression *AssignmentExpression) {
	astPrinter.write("AssignmentExpression")
	astPrinter.identationLevel++
	defer func() { astPrinter.identationLevel-- }()
	assignementExpression.Value.Accept(&astPrinter)
}

func (astPrinter AstPrinter) VisitBinaryExpression(binaryExpression *BinaryExpression) {
	astPrinter.write("BinaryExpression")
	astPrinter.identationLevel++
	binaryExpression.LHS.Accept(&astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel++
	binaryExpression.RHS.Accept(&astPrinter)
	astPrinter.identationLevel--
}

func (astPrinter AstPrinter) VisitCallExpression(callExpression *CallExpression) {
}

func (astPrinter AstPrinter) VisitGetExpression(getExpression *GetExpression) {
}

func (astPrinter AstPrinter) VisitGroupingExpression(groupingExpression *GroupingExpression) {
	astPrinter.write("GroupingExpression")
	astPrinter.identationLevel++
	defer func() { astPrinter.identationLevel-- }()
	groupingExpression.Expr.Accept(&astPrinter)
}

func (astPrinter AstPrinter) VisitLiteralExpression(literalExpression *LiteralExpression) {
	astPrinter.write(fmt.Sprintf("LiteralExpression (%s)", literalExpression.value))
}

func (astPrinter AstPrinter) VisitLogicalExpression(logicalExpression *LogicalExpression) {

}

func (astPrinter AstPrinter) VisitSetExpression(setExpression *SetExpression) {

}

func (astPrinter AstPrinter) VisitSuperExpression(superExpression *SuperExpression) {

}

func (astPrinter AstPrinter) VisitThisExpression(thisExpression *ThisExpression) {

}

func (astPrinter AstPrinter) VisitUnaryExpression(unaryExpression *UnaryExpression) {

}

func (astPrinter AstPrinter) VisitVariableExpression(variableExpression *VariableExpression) {

}

func (astPrinter *AstPrinter) VisitBlockStatement(blockStatement *BlockStatement) {
	astPrinter.identationLevel++
	defer func() { astPrinter.identationLevel-- }()
	for _, statement := range blockStatement.Statements {
		statement.Accept(astPrinter)
	}
}

func (astPrinter AstPrinter) VisitClassStatement(classStatement *ClassStatement) {

}

func (astPrinter AstPrinter) VisitExpressionStatement(expressionStatement *ExpressionStatement) {

}

func (astPrinter AstPrinter) VisitFunctionStatement(functionStatement *FunctionStatement) {

}

func (astPrinter AstPrinter) VisitIfStatement(ifStatement *IfStatement) {

}

func (astPrinter AstPrinter) VisitPrintStatement(printStatement *PrintStatement) {
	astPrinter.write("PrintStatement")
	astPrinter.identationLevel++
	defer func() { astPrinter.identationLevel-- }()
	printStatement.Expression.Accept(&astPrinter)
}

func (astPrinter AstPrinter) VisitReturnStatement(returnStatement *ReturnStatement) {
	astPrinter.write("ReturnStatement")
	astPrinter.identationLevel++
	defer func() { astPrinter.identationLevel-- }()
	returnStatement.Value.Accept(&astPrinter)
}

func (astPrinter AstPrinter) VisitVariableStatement(variableStatement *VariableStatement) {
	astPrinter.write("AssignmentExpression")
	astPrinter.identationLevel++
	astPrinter.write(" name: " + variableStatement.Name.Lexeme)
	astPrinter.write(" value: ")
	astPrinter.identationLevel++
	variableStatement.Initializer.Accept(&astPrinter)
	astPrinter.identationLevel--
	astPrinter.identationLevel--
}

func (astPrinter AstPrinter) VisitWhileStatement(whileStatement *WhileStatement) {

}

func (astPrinter AstPrinter) write(value string) {
	fmt.Fprintf(astPrinter.output, "%s%s\n", strings.Repeat(" ", astPrinter.identationLevel*astPrinter.tabSize), value)
}
