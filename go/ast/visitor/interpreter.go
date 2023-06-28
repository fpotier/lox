package visitor

import (
	"fmt"
	"reflect"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Interpreter struct {
	// Value can be virtually anything (string, number, boolean, object, nil, etc.)
	Value           ast.LoxValue
	environment     *ast.Environment
	HadRuntimeError bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Value:           nil,
		environment:     ast.NewEnvironment(),
		HadRuntimeError: false,
	}
}

func (visitor *Interpreter) Eval(statements []ast.Statement) {
	// This allows us to emulate a try/catch mechanism to exit the visitor as soon as possible
	// without changing the Visit...() methods to return an error and propagate manually these errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(loxerror.RuntimeError); ok {
				visitor.HadRuntimeError = true
				fmt.Println(err)
				return
			}
			panic(r)
		}
	}()

	for _, statement := range statements {
		visitor.execute(statement)
	}
}

func (visitor *Interpreter) VisitExpressionStatement(expressionStatement *ast.ExpressionStatement) {
	visitor.evaluate(expressionStatement.Expression)
}

func (visitor *Interpreter) VisitPrintStatement(printStatement *ast.PrintStatement) {
	value := visitor.evaluate(printStatement.Expression)
	if value == nil {
		fmt.Println("<nil>")
	} else {
		fmt.Println(value.String())
	}
}

func (visitor *Interpreter) VisitVariableStatement(variableStatement *ast.VariableStatement) {
	var value ast.LoxValue
	if variableStatement.Initializer != nil {
		value = visitor.evaluate(variableStatement.Initializer)
	}

	visitor.environment.Define(variableStatement.Name.Lexeme, value)
}

func (visitor *Interpreter) VisitBlockStatement(blockStatement *ast.BlockStatement) {
	visitor.executeBlock(blockStatement.Statements, ast.NewSubEnvironment(visitor.environment))
}

func (visitor *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := visitor.evaluate(binaryExpression.LHS)
	rhs := visitor.evaluate(binaryExpression.RHS)

	switch binaryExpression.Operator.Type {
	case lexer.Plus:
		switch {
		case lhs.IsNumber() && rhs.IsNumber():
			visitor.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value + rhs.(*ast.NumberValue).Value)
		case lhs.IsString() && rhs.IsString():
			visitor.Value = ast.NewStringValue(lhs.(*ast.StringValue).Value + rhs.(*ast.StringValue).Value)
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
		visitor.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value - rhs.(*ast.NumberValue).Value)
	case lexer.Star:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value * rhs.(*ast.NumberValue).Value)
	case lexer.Slash:
		// TODO check if rhs is 0 ?
		// go seems kinda broken: float64 / 0 = +Inf
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value / rhs.(*ast.NumberValue).Value)
	case lexer.Greater:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value > rhs.(*ast.NumberValue).Value)
	case lexer.GreaterEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value >= rhs.(*ast.NumberValue).Value)
	case lexer.Less:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value < rhs.(*ast.NumberValue).Value)
	case lexer.LessEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = ast.NewBooleanValue(lhs.(*ast.NumberValue).Value <= rhs.(*ast.NumberValue).Value)
	case lexer.EqualEqual:
		visitor.Value = ast.NewBooleanValue(lhs.Equals(rhs))
	case lexer.BangEqual:
		visitor.Value = ast.NewBooleanValue(!lhs.Equals(rhs))
	}
}

func (visitor *Interpreter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	visitor.Value = visitor.evaluate(groupingExpression.Expr)
}

func (visitor *Interpreter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	visitor.Value = literalExpression.LoxValue()
}

func (visitor *Interpreter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	rhs := visitor.evaluate(unaryExpression.RHS)

	switch unaryExpression.Operator.Type {
	case lexer.Bang:
		visitor.Value = ast.NewBooleanValue(!rhs.IsTruthy())
	case lexer.Dash:
		assertNumberOperand(unaryExpression.Operator, rhs)
		visitor.Value = ast.NewNumberValue(-rhs.(*ast.NumberValue).Value)
	}
}

func (visitor *Interpreter) VisitVariableExpression(variableExpression *ast.VariableExpression) {
	visitor.Value = visitor.environment.Get(variableExpression.Name)
}

func (visitor *Interpreter) VisitAssignmentExpression(assignmentExpression *ast.AssignmentExpression) {
	value := visitor.evaluate(assignmentExpression.Value)
	visitor.environment.Assign(assignmentExpression.Name, value)
	visitor.Value = value
}

func (visitor *Interpreter) executeBlock(statements []ast.Statement, subEnvironment *ast.Environment) {
	previousEnv := visitor.environment
	visitor.environment = subEnvironment
	// TODO error handling
	for _, statement := range statements {
		visitor.execute(statement)
	}

	visitor.environment = previousEnv
}

func (visitor *Interpreter) execute(statement ast.Statement) {
	statement.Accept(visitor)
}

func (visitor *Interpreter) evaluate(expression ast.Expression) ast.LoxValue {
	// TODO: handle error -> causes nil pointer dereference
	newVisitor := &Interpreter{
		Value:           nil,
		environment:     visitor.environment,
		HadRuntimeError: false,
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
