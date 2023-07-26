package ast

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Interpreter struct {
	// Value can be virtually anything (string, number, boolean, object, nil, etc.)
	Value           LoxValue
	HadRuntimeError bool
	OutputStream    io.Writer
	globals         *Environment
	environment     *Environment
	locals          map[Expression]int
}

func NewInterpreter() *Interpreter {
	i := Interpreter{
		Value:           nil,
		HadRuntimeError: false,
		OutputStream:    os.Stdout,
		globals:         NewEnvironment(),
		locals:          make(map[Expression]int),
	}
	i.environment = i.globals

	nativeClock := NativeFunction{
		name:  "clock",
		arity: 0,
		code: func(*Interpreter, []LoxValue) LoxValue {
			return NewNumberValue(float64(time.Now().UnixMilli()) / 1000.0)
		},
	}

	i.globals.Define(nativeClock.name, nativeClock)

	return &i
}

func (i *Interpreter) Eval(statements []Statement) {
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

func (i *Interpreter) VisitIfStatement(ifStatment *IfStatement) {
	if i.evaluate(ifStatment.Condition).IsTruthy() {
		i.execute(ifStatment.ThenCode)
	} else if ifStatment.ElseCode != nil {
		i.execute(ifStatment.ElseCode)
	}
}

func (i *Interpreter) VisitWhileStatement(whileStatement *WhileStatement) {
	for i.evaluate(whileStatement.Condition).IsTruthy() {
		i.execute(whileStatement.Body)
	}
}

func (i *Interpreter) VisitExpressionStatement(expressionStatement *ExpressionStatement) {
	i.evaluate(expressionStatement.Expression)
}

func (i *Interpreter) VisitPrintStatement(printStatement *PrintStatement) {
	value := i.evaluate(printStatement.Expression)
	if value == nil {
		fmt.Fprint(i.OutputStream, "<nil>\n")
	} else {
		fmt.Fprintf(i.OutputStream, "%s\n", value.String())
	}
}

func (i *Interpreter) VisitVariableStatement(variableStatement *VariableStatement) {
	var value LoxValue
	if variableStatement.Initializer != nil {
		value = i.evaluate(variableStatement.Initializer)
	}

	i.environment.Define(variableStatement.Name.Lexeme, value)
}

func (i *Interpreter) VisitBlockStatement(blockStatement *BlockStatement) {
	i.executeBlock(blockStatement.Statements, NewSubEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStatement(classStatement *ClassStatement) {
	i.environment.Define(classStatement.Name.Lexeme, nil)
	class := NewLoxClass(classStatement.Name.Lexeme)
	i.environment.Assign(classStatement.Name, class)
}

func (i *Interpreter) VisitFunctionStatement(functionStatement *FunctionStatement) {
	function := NewLoxFunction(functionStatement, i.environment)
	i.environment.Define(functionStatement.Name.Lexeme, function)
}

func (i *Interpreter) VisitReturnStatement(returnStatement *ReturnStatement) {
	var value LoxValue
	if returnStatement.Value != nil {
		value = i.evaluate(returnStatement.Value)
	}

	panic(value)
}

func (i *Interpreter) VisitBinaryExpression(binaryExpression *BinaryExpression) {
	lhs := i.evaluate(binaryExpression.LHS)
	rhs := i.evaluate(binaryExpression.RHS)

	switch binaryExpression.Operator.Type {
	case lexer.Plus:
		switch {
		case lhs.IsNumber() && rhs.IsNumber():
			i.Value = NewNumberValue(lhs.(*NumberValue).Value + rhs.(*NumberValue).Value)
		case lhs.IsString() && rhs.IsString():
			i.Value = NewStringValue(lhs.(*StringValue).Value + rhs.(*StringValue).Value)
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
		i.Value = NewNumberValue(lhs.(*NumberValue).Value - rhs.(*NumberValue).Value)
	case lexer.Star:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewNumberValue(lhs.(*NumberValue).Value * rhs.(*NumberValue).Value)
	case lexer.Slash:
		// TODO check if rhs is 0 ?
		// go seems kinda broken: float64 / 0 = +Inf
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewNumberValue(lhs.(*NumberValue).Value / rhs.(*NumberValue).Value)
	case lexer.Greater:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewBooleanValue(lhs.(*NumberValue).Value > rhs.(*NumberValue).Value)
	case lexer.GreaterEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewBooleanValue(lhs.(*NumberValue).Value >= rhs.(*NumberValue).Value)
	case lexer.Less:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewBooleanValue(lhs.(*NumberValue).Value < rhs.(*NumberValue).Value)
	case lexer.LessEqual:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = NewBooleanValue(lhs.(*NumberValue).Value <= rhs.(*NumberValue).Value)
	case lexer.EqualEqual:
		i.Value = NewBooleanValue(lhs.Equals(rhs))
	case lexer.BangEqual:
		i.Value = NewBooleanValue(!lhs.Equals(rhs))
	}
}

func (i *Interpreter) VisitLogicalExpression(logicalExpression *LogicalExpression) {
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

func (i *Interpreter) VisitGroupingExpression(groupingExpression *GroupingExpression) {
	i.Value = i.evaluate(groupingExpression.Expr)
}

func (i *Interpreter) VisitLiteralExpression(literalExpression *LiteralExpression) {
	i.Value = literalExpression.LoxValue()
}

func (i *Interpreter) VisitUnaryExpression(unaryExpression *UnaryExpression) {
	rhs := i.evaluate(unaryExpression.RHS)

	switch unaryExpression.Operator.Type {
	case lexer.Bang:
		i.Value = NewBooleanValue(!rhs.IsTruthy())
	case lexer.Dash:
		assertNumberOperand(unaryExpression.Operator, rhs)
		i.Value = NewNumberValue(-rhs.(*NumberValue).Value)
	}
}

func (i *Interpreter) VisitVariableExpression(variableExpression *VariableExpression) {
	i.Value = i.lookupVariable(variableExpression.Name, variableExpression)
}

func (i *Interpreter) VisitAssignmentExpression(assignmentExpression *AssignmentExpression) {
	value := i.evaluate(assignmentExpression.Value)
	if distance, ok := i.locals[assignmentExpression]; ok {
		i.environment.AssignAt(distance, assignmentExpression.Name, value)
	} else {
		i.globals.Assign(assignmentExpression.Name, value)
	}
	i.Value = value
}

func (i *Interpreter) VisitCallExpression(callExpression *CallExpression) {
	callee := i.evaluate(callExpression.Callee)

	arguments := make([]LoxValue, 0)
	for _, argument := range callExpression.Args {
		arguments = append(arguments, i.evaluate(argument))
	}

	if function, ok := callee.(LoxCallable); ok {
		if function.Arity() != len(arguments) {
			panic(loxerror.RuntimeError{
				Message: fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(arguments)),
			})
		}
		i.Value = function.Call(i, arguments)
	} else {
		panic(loxerror.RuntimeError{Message: "Can only call functions and classes"})
	}
}

func (i *Interpreter) VisitGetExpression(getExpression *GetExpression) {
	object := i.evaluate(getExpression.Object)
	if object, ok := object.(*LoxInstance); ok {
		i.Value = object.Get(getExpression.Name)
		return
	}

	panic(loxerror.RuntimeError{Message: "Only instances have properties"})
}

func (i *Interpreter) VisitSetExpression(setExpression *SetExpression) {
	object := i.evaluate(setExpression.Object)

	if object, ok := object.(*LoxInstance); ok {
		value := i.evaluate(setExpression.Value)
		object.Set(setExpression.Name, value)
		i.Value = value

		return
	}
	panic(loxerror.RuntimeError{Message: "Only instances have fields"})
}

func (i *Interpreter) executeBlock(statements []Statement, subEnvironment *Environment) {
	previousEnv := i.environment
	i.environment = subEnvironment
	// TODO error handling
	for _, statement := range statements {
		i.execute(statement)
	}

	i.environment = previousEnv
}

func (i *Interpreter) execute(statement Statement) {
	statement.Accept(i)
}

func (i *Interpreter) evaluate(expression Expression) LoxValue {
	// TODO: handle error -> causes nil pointer dereference
	newVisitor := NewInterpreter()
	newVisitor.environment = i.environment
	newVisitor.globals = i.globals
	newVisitor.locals = i.locals

	expression.Accept(newVisitor)

	return newVisitor.Value
}

func (i *Interpreter) resolve(e Expression, depth int) {
	i.locals[e] = depth
}

func (i *Interpreter) lookupVariable(name lexer.Token, e Expression) LoxValue {
	if distance, ok := i.locals[e]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func assertNumberOperands(operator lexer.Token, lhs LoxValue, rhs LoxValue) {
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

func assertNumberOperand(operator lexer.Token, rhs LoxValue) {
	if !rhs.IsNumber() {
		panic(loxerror.RuntimeError{
			// TODO: print lox types instead of go types
			Message: fmt.Sprintf("Operator '%v': incompatible type %v", operator.Lexeme, reflect.TypeOf(rhs)),
		})
	}
}
