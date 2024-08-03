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

func (astPrinter *AstPrinter) VisitAssignmentExpression(assignementExpression *AssignmentExpression) {
	astPrinter.write("AssignmentExpression")
	astPrinter.identationLevel++

	astPrinter.write("identifier: " + assignementExpression.Name.Lexeme)

	if assignementExpression.Value != nil {
		astPrinter.write("value: ")
		astPrinter.identationLevel++
		assignementExpression.Value.Accept(astPrinter)
		astPrinter.identationLevel--
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitBinaryExpression(binaryExpression *BinaryExpression) {
	astPrinter.write("BinaryExpression")
	astPrinter.identationLevel++

	astPrinter.write("operator: " + binaryExpression.Operator.Lexeme)

	astPrinter.write("left_operand: ")
	astPrinter.identationLevel++
	binaryExpression.LHS.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("right_operand: ")
	astPrinter.identationLevel++
	binaryExpression.RHS.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitCallExpression(callExpression *CallExpression) {
	astPrinter.write("CallExpression")
	astPrinter.identationLevel++

	astPrinter.write("callee: ")
	astPrinter.identationLevel++
	callExpression.Callee.Accept(astPrinter)
	astPrinter.identationLevel--

	if len(callExpression.Args) > 0 {
		astPrinter.writeExpression("arguments", callExpression.Args)
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitGetExpression(getExpression *GetExpression) {
	astPrinter.write("GetExpression")
	astPrinter.identationLevel++

	astPrinter.write("object:")
	astPrinter.identationLevel++
	getExpression.Object.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("property: " + getExpression.Name.Lexeme)

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitGroupingExpression(groupingExpression *GroupingExpression) {
	astPrinter.write("GroupingExpression")
	astPrinter.identationLevel++

	astPrinter.write("expression")
	astPrinter.identationLevel++
	groupingExpression.Expr.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitLiteralExpression(literalExpression *LiteralExpression) {
	astPrinter.write("LiteralExpression")
	astPrinter.identationLevel++
	astPrinter.write("value: " + literalExpression.value.String())
	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitLogicalExpression(logicalExpression *LogicalExpression) {
	astPrinter.write("LogicalExpression")
	astPrinter.identationLevel++

	astPrinter.write("operator: " + logicalExpression.Operator.Lexeme)

	astPrinter.write("left_operand: ")
	astPrinter.identationLevel++
	logicalExpression.LHS.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("right_operand: ")
	astPrinter.identationLevel++
	logicalExpression.RHS.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitSetExpression(setExpression *SetExpression) {
	astPrinter.write("SetExpression")
	astPrinter.identationLevel++

	astPrinter.write("object:")
	astPrinter.identationLevel++
	setExpression.Object.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("property: " + setExpression.Name.Lexeme)

	astPrinter.write("value:")
	astPrinter.identationLevel++
	setExpression.Value.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitSuperExpression(superExpression *SuperExpression) {
	astPrinter.write("SuperExpression")
	astPrinter.identationLevel++
	astPrinter.write("method: " + superExpression.Method.Lexeme)
	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitThisExpression(thisExpression *ThisExpression) {
	astPrinter.write("ThisExpression")
}

func (astPrinter *AstPrinter) VisitUnaryExpression(unaryExpression *UnaryExpression) {
	astPrinter.write("UnaryExpression")
	astPrinter.identationLevel++

	astPrinter.write("operator: " + unaryExpression.Operator.Lexeme)

	astPrinter.write("operand: ")
	astPrinter.identationLevel++
	unaryExpression.RHS.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitVariableExpression(variableExpression *VariableExpression) {
	astPrinter.write("VariableExpression")
	astPrinter.identationLevel++
	astPrinter.write("name: " + variableExpression.Name.Lexeme)
	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitBlockStatement(blockStatement *BlockStatement) {
	astPrinter.write("BlockStatement")
	astPrinter.identationLevel++

	astPrinter.writeStatements("statements", blockStatement.Statements)

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitClassStatement(classStatement *ClassStatement) {
	astPrinter.write("ClassStatement")
	astPrinter.identationLevel++

	astPrinter.write("name: " + classStatement.Name.Lexeme)
	if classStatement.Superclass != nil {
		astPrinter.write("superclass:")
		astPrinter.identationLevel++
		classStatement.Superclass.Accept(astPrinter)
		astPrinter.identationLevel--
	}

	if len(classStatement.Methods) > 0 {
		// FIXME: why are we using pointers?
		// astPrinter.writeStatements("methods", classStatement.Methods)
		astPrinter.write("methods:")
		astPrinter.identationLevel++
		for i, method := range classStatement.Methods {
			astPrinter.identationLevel++
			astPrinter.write(fmt.Sprintf("%d: ", i))

			astPrinter.identationLevel++
			method.Accept(astPrinter)
			astPrinter.identationLevel--

			astPrinter.identationLevel--

		}
		astPrinter.identationLevel--
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitExpressionStatement(expressionStatement *ExpressionStatement) {
	astPrinter.write("ExpressionStatement")
	astPrinter.identationLevel++

	astPrinter.write("expression:")
	astPrinter.identationLevel++
	expressionStatement.Expression.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitFunctionStatement(functionStatement *FunctionStatement) {
	astPrinter.write("FunctionStatement")
	astPrinter.identationLevel++

	astPrinter.write("name: " + functionStatement.Name.Lexeme)

	if len(functionStatement.Parameters) > 0 {
		for i, parameter := range functionStatement.Parameters {
			astPrinter.identationLevel++
			astPrinter.write(fmt.Sprintf("%d: %s", i, parameter.Lexeme))
			astPrinter.identationLevel--
		}
	}

	astPrinter.writeStatements("body", functionStatement.Body)

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitIfStatement(ifStatement *IfStatement) {
	astPrinter.write("IfStatement")
	astPrinter.identationLevel++

	astPrinter.write("condition:")
	astPrinter.identationLevel++
	ifStatement.Condition.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("then:")
	astPrinter.identationLevel++
	ifStatement.ThenCode.Accept(astPrinter)
	astPrinter.identationLevel--

	if ifStatement.ElseCode != nil {
		astPrinter.write("else:")
		astPrinter.identationLevel++
		ifStatement.ElseCode.Accept(astPrinter)
		astPrinter.identationLevel--
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitPrintStatement(printStatement *PrintStatement) {
	astPrinter.write("PrintStatement")
	astPrinter.identationLevel++

	astPrinter.write("value:")
	astPrinter.identationLevel++
	printStatement.Expression.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitReturnStatement(returnStatement *ReturnStatement) {
	astPrinter.write("ReturnStatement")
	astPrinter.identationLevel++

	if returnStatement.Value != nil {
		astPrinter.write("value:")
		astPrinter.identationLevel++
		returnStatement.Value.Accept(astPrinter)
		astPrinter.identationLevel--
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitVariableStatement(variableStatement *VariableStatement) {
	astPrinter.write("VariableStatement")
	astPrinter.identationLevel++

	astPrinter.write("name: " + variableStatement.Name.Lexeme)

	if variableStatement.Initializer != nil {
		astPrinter.write("initializer: ")
		astPrinter.identationLevel++
		variableStatement.Initializer.Accept(astPrinter)
		astPrinter.identationLevel--
	}

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) VisitWhileStatement(whileStatement *WhileStatement) {
	astPrinter.write("WhileStatement")
	astPrinter.identationLevel++

	astPrinter.write("condition:")
	astPrinter.identationLevel++
	whileStatement.Condition.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.write("body:")
	astPrinter.identationLevel++
	whileStatement.Body.Accept(astPrinter)
	astPrinter.identationLevel--

	astPrinter.identationLevel--
}

func (astPrinter *AstPrinter) write(value string) {
	fmt.Fprintf(astPrinter.output,
		"%s%s\n",
		strings.Repeat(" ", astPrinter.identationLevel*astPrinter.tabSize),
		value,
	)
}

func (astPrinter *AstPrinter) writeStatements(name string, statements []Statement) {
	astPrinter.write(name + ": ")
	for i, statement := range statements {
		astPrinter.identationLevel++
		astPrinter.write(fmt.Sprintf("%d: ", i))

		astPrinter.identationLevel++
		statement.Accept(astPrinter)
		astPrinter.identationLevel--

		astPrinter.identationLevel--
	}
}

func (astPrinter *AstPrinter) writeExpression(name string, expressions []Expression) {
	astPrinter.write(name + ": ")
	for i, expression := range expressions {
		astPrinter.identationLevel++
		astPrinter.write(fmt.Sprintf("%d: ", i))

		astPrinter.identationLevel++
		expression.Accept(astPrinter)
		astPrinter.identationLevel--

		astPrinter.identationLevel--
	}
}