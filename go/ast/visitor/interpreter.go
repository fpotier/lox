package visitor

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Interpreter struct {
	// Value can be virtually anything (string, number, boolean, object, nil, etc.)
	Value           ast.LoxValue
	HadRuntimeError bool
	OutputStream    io.Writer
	environment     *ast.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Value:           nil,
		HadRuntimeError: false,
		OutputStream:    os.Stdout,
		environment:     ast.NewEnvironment(),
	}
}

func (i *Interpreter) Eval(statements []ast.Statement) {
	// This allows us to emulate a try/catch mechanism to exit the visitor as soon as possible
	// without changing the Visit...() methods to return an error and propagate manually these errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(loxerror.RuntimeError); ok {
				i.HadRuntimeError = true
				fmt.Println(err)
				return
			}
			panic(r)
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitIfStatement(ifStatment *ast.IfStatement) {
	if i.evaluate(ifStatment.Condition).IsTruthy() {
		i.execute(ifStatment.ThenCode)
	} else if ifStatment.ElseCode != nil {
		i.execute(ifStatment.ElseCode)
	}
}

func (i *Interpreter) VisitWhileStatement(whileStatement *ast.WhileStatement) {
	for i.evaluate(whileStatement.Condition).IsTruthy() {
		i.execute(whileStatement.Body)
	}
}

func (i *Interpreter) VisitExpressionStatement(expressionStatement *ast.ExpressionStatement) {
	i.evaluate(expressionStatement.Expression)
}

func (i *Interpreter) VisitPrintStatement(printStatement *ast.PrintStatement) {
	value := i.evaluate(printStatement.Expression)
	if value == nil {
		fmt.Fprint(i.OutputStream, "<nil>\n")
	} else {
		fmt.Fprintf(i.OutputStream, "%s\n", value.String())
	}
}

func (i *Interpreter) VisitVariableStatement(variableStatement *ast.VariableStatement) {
	var value ast.LoxValue
	if variableStatement.Initializer != nil {
		value = i.evaluate(variableStatement.Initializer)
	}

	i.environment.Define(variableStatement.Name.Lexeme, value)
}

func (i *Interpreter) VisitBlockStatement(blockStatement *ast.BlockStatement) {
	i.executeBlock(blockStatement.Statements, ast.NewSubEnvironment(i.environment))
}

func (i *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := i.evaluate(binaryExpression.LHS)
	rhs := i.evaluate(binaryExpression.RHS)

	switch binaryExpression.Operator.Type {
	case lexer.Plus:
		switch {
		case lhs.IsNumber() && rhs.IsNumber():
			i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value + rhs.(*ast.NumberValue).Value)
		case lhs.IsString() && rhs.IsString():
			i.Value = ast.NewStringValue(lhs.(*ast.StringValue).Value + rhs.(*ast.StringValue).Value)
		default:
			// TODO: print lox types instead of go types
			panic(loxerror.RuntimeError{
				Message: fmt.Sprintf("Operator '%v': incompatible types %v and %v",
					binaryExpression.Operator.Lexeme,
					reflect.TypeOf(lhs),
					reflect.TypeOf(rhs)),
			})
		}
	case lexer.Dash:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value - rhs.(*ast.NumberValue).Value)
	case lexer.Star:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value * rhs.(*ast.NumberValue).Value)
	case lexer.Slash:
		// TODO check if rhs is 0 ?
		// go seems kinda broken: float64 / 0 = +Inf
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value / rhs.(*ast.NumberValue).Value)
	case lexer.Greater:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value > rhs.(*ast.NumberValue).Value)
	case lexer.GreaterEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value >= rhs.(*ast.NumberValue).Value)
	case lexer.Less:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value < rhs.(*ast.NumberValue).Value)
	case lexer.LessEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value <= rhs.(*ast.NumberValue).Value)
	case lexer.EqualEqual:
		i.Value = ast.NewBooleanValue(lhs.Equals(rhs))
	case lexer.BangEqual:
		i.Value = ast.NewBooleanValue(!lhs.Equals(rhs))
	}
}

func (i *Interpreter) VisitLogicalExpression(logicalExpression *ast.LogicalExpression) {
	lhs := i.evaluate(logicalExpression.LHS)
	switch {
	case logicalExpression.Operator.Type == lexer.Or && lhs.IsTruthy():
		i.Value = lhs
	case logicalExpression.Operator.Type == lexer.And && !lhs.IsTruthy():
		i.Value = lhs
	default:
		i.Value = i.evaluate(logicalExpression.RHS)
	}
}

func (i *Interpreter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	i.Value = i.evaluate(groupingExpression.Expr)
}

func (i *Interpreter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	i.Value = literalExpression.LoxValue()
}

func (i *Interpreter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	rhs := i.evaluate(unaryExpression.RHS)

	switch unaryExpression.Operator.Type {
	case lexer.Bang:
		i.Value = ast.NewBooleanValue(!rhs.IsTruthy())
	case lexer.Dash:
		assertNumberOperand(unaryExpression.Operator, rhs)
		i.Value = ast.NewNumberValue(-rhs.(*ast.NumberValue).Value)
	}
}

func (i *Interpreter) VisitVariableExpression(variableExpression *ast.VariableExpression) {
	i.Value = i.environment.Get(variableExpression.Name)
}

func (i *Interpreter) VisitAssignmentExpression(assignmentExpression *ast.AssignmentExpression) {
	value := i.evaluate(assignmentExpression.Value)
	i.environment.Assign(assignmentExpression.Name, value)
	i.Value = value
}

func (i *Interpreter) executeBlock(statements []ast.Statement, subEnvironment *ast.Environment) {
	previousEnv := i.environment
	i.environment = subEnvironment
	// TODO error handling
	for _, statement := range statements {
		i.execute(statement)
	}

	i.environment = previousEnv
}

func (i *Interpreter) execute(statement ast.Statement) {
	statement.Accept(i)
}

func (i *Interpreter) evaluate(expression ast.Expression) ast.LoxValue {
	// TODO: handle error -> causes nil pointer dereference
	newVisitor := &Interpreter{
		Value:           nil,
		HadRuntimeError: false,
		OutputStream:    os.Stdout,
		environment:     i.environment,
	}
	expression.Accept(newVisitor)

	return newVisitor.Value
}

func assertNumberOperands(operator lexer.Token, lhs ast.LoxValue, rhs ast.LoxValue) {
	if !lhs.IsNumber() || !rhs.IsNumber() {
		panic(loxerror.RuntimeError{
			// TODO: print lox types instead of go types
			Message: fmt.Sprintf("Operator '%v': incompatible types %v and %v",
				operator.Lexeme,
				reflect.TypeOf(lhs),
				reflect.TypeOf(rhs)),
		})
	}
}

func assertNumberOperand(operator lexer.Token, rhs ast.LoxValue) {
	if !rhs.IsNumber() {
		panic(loxerror.RuntimeError{
			// TODO: print lox types instead of go types
			Message: fmt.Sprintf("Operator '%v': incompatible type %v", operator.Lexeme, reflect.TypeOf(rhs)),
		})
	}
}