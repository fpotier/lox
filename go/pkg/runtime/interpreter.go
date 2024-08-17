//go:generate go run ../../cmd/code-generator runtime

package runtime

import (
	"fmt"
	"io"

	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/lexer"
	"github.com/fpotier/lox/go/pkg/loxerror"
)

type Interpreter struct {
	// Value can be virtually anything (string, number, boolean, object, nil, etc.)
	Value           ast.LoxValue
	HadRuntimeError bool
	ErrorFormatter  loxerror.ErrorFormatter
	OutputStream    io.Writer
	globals         *Environment
	environment     *Environment
	locals          map[ast.Expression]int
}

func NewInterpreter(outputStream io.Writer, errorFormatter loxerror.ErrorFormatter) *Interpreter {
	i := Interpreter{
		Value:           ast.NewNilValue(),
		HadRuntimeError: false,
		ErrorFormatter:  errorFormatter,
		OutputStream:    outputStream,
		globals:         NewEnvironment(),
		environment:     nil,
		locals:          make(map[ast.Expression]int),
	}
	i.environment = i.globals

	for _, nativeFunction := range builtinNativeFunctions {
		i.globals.Define(nativeFunction.name, nativeFunction)
	}

	return &i
}

func (i *Interpreter) Eval(statements []ast.Statement) {
	// This allows us to emulate a try/catch mechanism to exit the visitor as soon as possible
	// without changing the Visit...() methods to return an error and propagate manually these errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(loxerror.LoxError); ok {
				i.HadRuntimeError = true
				// TODO: better runtime error messages
				i.ErrorFormatter.PushError(err)
				return
			}
			panic(r)
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitAssignmentExpression(assignmentExpression *ast.AssignmentExpression) {
	value := i.evaluate(assignmentExpression.Value)
	if distance, ok := i.locals[assignmentExpression]; ok {
		i.environment.AssignAt(distance, assignmentExpression.Name, value)
	} else {
		i.globals.Assign(assignmentExpression.Name, value)
	}
	i.Value = value
}

func (i *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := i.evaluate(binaryExpression.LHS)
	rhs := i.evaluate(binaryExpression.RHS)

	switch binaryExpression.Operator.Type {
	case lexer.Plus:
		switch {
		case lhs.Kind() == ast.Number && rhs.Kind() == ast.Number:
			i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value + rhs.(*ast.NumberValue).Value)
		case lhs.Kind() == ast.String && rhs.Kind() == ast.String:
			i.Value = ast.NewStringValue(lhs.(*ast.StringValue).Value + rhs.(*ast.StringValue).Value)
		default:
			panic(NewUnsupportedBinaryOperation(binaryExpression.Operator.Line,
				binaryExpression.Operator.Lexeme,
				ast.KindString[lhs.Kind()],
				ast.KindString[rhs.Kind()]))
		}
	case lexer.Dash:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value - rhs.(*ast.NumberValue).Value)
	case lexer.Star:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		i.Value = ast.NewNumberValue(lhs.(*ast.NumberValue).Value * rhs.(*ast.NumberValue).Value)
	case lexer.Slash:
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

func (i *Interpreter) VisitCallExpression(callExpression *ast.CallExpression) {
	callee := i.evaluate(callExpression.Callee)

	arguments := make([]ast.LoxValue, 0)
	for _, argument := range callExpression.Args {
		arguments = append(arguments, i.evaluate(argument))
	}

	if function, ok := callee.(LoxCallable); ok {
		if function.Arity() != len(arguments) {
			panic(NewBadArity(callExpression.Position.Line, function.Name(), function.Arity(), len(arguments)))
		}
		i.Value = function.Call(i, arguments)
	} else {
		panic(NewNotCallable(callExpression.Position.Line))
	}
}

func (i *Interpreter) VisitGetExpression(getExpression *ast.GetExpression) {
	object := i.evaluate(getExpression.Object)
	if object, ok := object.(*LoxInstance); ok {
		i.Value = object.Get(getExpression.Name)
		return
	}

	panic(NewInvalidSetGet(getExpression.Name.Line))
}

func (i *Interpreter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	i.Value = i.evaluate(groupingExpression.Expr)
}

func (i *Interpreter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	i.Value = literalExpression.LoxValue()
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

func (i *Interpreter) VisitSetExpression(setExpression *ast.SetExpression) {
	object := i.evaluate(setExpression.Object)

	if object, ok := object.(*LoxInstance); ok {
		value := i.evaluate(setExpression.Value)
		object.Set(setExpression.Name, value)
		i.Value = value

		return
	}

	panic(NewInvalidSetGet(setExpression.Name.Line))
}

func (i *Interpreter) VisitSuperExpression(superExpression *ast.SuperExpression) {
	distance := i.locals[superExpression]
	superclass := i.environment.GetAt(distance, "super").(*LoxClass)
	this := i.environment.GetAt(distance-1, "this").(*LoxInstance)
	method, ok := superclass.findMethod(superExpression.Method.Lexeme)
	if !ok {
		panic(NewUndefinedProperty(superExpression.Method.Line, superExpression.Keyword.Lexeme, superclass.String()))
	}

	i.Value = method.Bind(this)
}

func (i *Interpreter) VisitThisExpression(thisExpression *ast.ThisExpression) {
	i.Value = i.lookupVariable(thisExpression.Keyword, thisExpression)
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
	i.Value = i.lookupVariable(variableExpression.Name, variableExpression)
}

func (i *Interpreter) VisitBlockStatement(blockStatement *ast.BlockStatement) {
	i.executeBlock(blockStatement.Statements, NewSubEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStatement(classStatement *ast.ClassStatement) {
	var superclass *LoxClass
	if classStatement.Superclass != nil {
		result := i.evaluate(classStatement.Superclass)
		var ok bool
		superclass, ok = result.(*LoxClass)
		if !ok {
			panic(NewInvalidInheritance(classStatement.Superclass.Name.Line, "Superclass must be a class"))
		}
	}

	i.environment.Define(classStatement.Name.Lexeme, ast.NewNilValue())

	if classStatement.Superclass != nil {
		i.environment = NewSubEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range classStatement.Methods {
		function := NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		function.setClassName(classStatement.Name.Lexeme)
		methods[method.Name.Lexeme] = function
	}

	class := NewLoxClass(classStatement.Name.Lexeme, superclass, methods)

	if superclass != nil {
		i.environment = i.environment.enclosing
	}

	i.environment.Assign(classStatement.Name, class)
}

func (i *Interpreter) VisitExpressionStatement(expressionStatement *ast.ExpressionStatement) {
	i.evaluate(expressionStatement.Expression)
}

func (i *Interpreter) VisitFunctionStatement(functionStatement *ast.FunctionStatement) {
	function := NewLoxFunction(functionStatement, i.environment, false)
	i.environment.Define(functionStatement.Name.Lexeme, function)
}

func (i *Interpreter) VisitIfStatement(ifStatment *ast.IfStatement) {
	if i.evaluate(ifStatment.Condition).IsTruthy() {
		i.execute(ifStatment.ThenCode)
	} else if ifStatment.ElseCode != nil {
		i.execute(ifStatment.ElseCode)
	}
}

func (i *Interpreter) VisitPrintStatement(printStatement *ast.PrintStatement) {
	value := i.evaluate(printStatement.Expression)
	fmt.Fprintf(i.OutputStream, "%s\n", value.String())
}

func (i *Interpreter) VisitReturnStatement(returnStatement *ast.ReturnStatement) {
	var value ast.LoxValue = ast.NewNilValue()
	if returnStatement.Value != nil {
		value = i.evaluate(returnStatement.Value)
	}

	panic(value)
}

func (i *Interpreter) VisitVariableStatement(variableStatement *ast.VariableStatement) {
	var value ast.LoxValue = ast.NewNilValue()
	if variableStatement.Initializer != nil {
		value = i.evaluate(variableStatement.Initializer)
	}

	i.environment.Define(variableStatement.Name.Lexeme, value)
}

func (i *Interpreter) VisitWhileStatement(whileStatement *ast.WhileStatement) {
	for i.evaluate(whileStatement.Condition).IsTruthy() {
		i.execute(whileStatement.Body)
	}
}

func (i *Interpreter) executeBlock(statements []ast.Statement, subEnvironment *Environment) {
	previousEnv := i.environment
	i.environment = subEnvironment
	defer func() { i.environment = previousEnv }()
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
	oldValue := i.Value
	expression.Accept(i)
	evalValue := i.Value
	i.Value = oldValue

	return evalValue
}

func (i *Interpreter) resolve(e ast.Expression, depth int) {
	i.locals[e] = depth
}

func (i *Interpreter) lookupVariable(name lexer.Token, e ast.Expression) ast.LoxValue {
	if distance, ok := i.locals[e]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}

	return i.globals.Get(name)
}

func assertNumberOperands(operator lexer.Token, lhs ast.LoxValue, rhs ast.LoxValue) {
	if !(lhs.Kind() == ast.Number) || !(rhs.Kind() == ast.Number) {
		panic(
			NewUnsupportedBinaryOperation(
				operator.Line,
				operator.Lexeme,
				ast.KindString[lhs.Kind()],
				ast.KindString[rhs.Kind()],
			),
		)
	}
}

func assertNumberOperand(operator lexer.Token, rhs ast.LoxValue) {
	if !(rhs.Kind() == ast.Number) {
		panic(NewUnsupportedUnaryOperation(operator.Line, operator.Lexeme, ast.KindString[rhs.Kind()]))
	}
}
