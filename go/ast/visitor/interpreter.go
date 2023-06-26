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
	Value    ast.LoxValue
	HadError bool
}

func (visitor *Interpreter) Eval(expression ast.Expression) interface{} {
	// This allows us to emulate a try/catch mechanism to exit the visitor as soon as possible
	// without changing the Visit...() methods to return an error and propagate manually these errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(loxerror.RuntimeError); ok {
				visitor.HadError = true
				fmt.Println(err)
				return
			}
			panic(r)
		}
	}()
	expression.Accept(visitor)
	return visitor.Value
}

func (visitor *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := visitor.evaluate(binaryExpression.Lhs)
	rhs := visitor.evaluate(binaryExpression.Rhs)

	switch binaryExpression.Operator.Type {
	case lexer.PLUS:
		if isNumber(lhs) && isNumber(rhs) {
			visitor.Value = &ast.NumberValue{
				Value: lhs.(*ast.NumberValue).Value + rhs.(*ast.NumberValue).Value,
			}
		} else if isString(lhs) && isString(rhs) {
			visitor.Value = &ast.StringValue{
				Value: lhs.(*ast.StringValue).Value + rhs.(*ast.StringValue).Value,
			}
		} else {
			// TODO: print lox types instead of go types
			panic(loxerror.RuntimeError{
				Message: fmt.Sprintf("Operator '%v': incompatible types %v and %v", binaryExpression.Operator.Lexeme, reflect.TypeOf(lhs), reflect.TypeOf(rhs)),
			})
		}
	case lexer.DASH:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value - rhs.(*ast.NumberValue).Value}
	case lexer.STAR:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value * rhs.(*ast.NumberValue).Value}
	case lexer.SLASH:
		assertNumberOperands(binaryExpression.Operator, lhs, rhs)
		visitor.Value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value / rhs.(*ast.NumberValue).Value}
	}
}

func (visitor *Interpreter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	visitor.Value = visitor.evaluate(groupingExpression.Expr)
}

func (visitor *Interpreter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	visitor.Value = literalExpression.LoxValue()
}

func (visitor *Interpreter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	rhs := visitor.evaluate(unaryExpression.Rhs)

	switch unaryExpression.Operator.Type {
	case lexer.BANG:
		visitor.Value = &ast.BooleanValue{Value: !rhs.IsTruthy()}
	case lexer.DASH:
		assertNumberOperand(unaryExpression.Operator, rhs)
		visitor.Value = &ast.NumberValue{Value: -rhs.(*ast.NumberValue).Value}
	}
}

func (visitor *Interpreter) evaluate(expression ast.Expression) ast.LoxValue {
	newVisitor := &Interpreter{}
	expression.Accept(newVisitor)

	return newVisitor.Value
}

func isNumber(value ast.LoxValue) bool {
	_, ok := value.(*ast.NumberValue)
	return ok
}

func isString(value ast.LoxValue) bool {
	_, ok := value.(*ast.StringValue)
	return ok
}

func assertNumberOperands(operator lexer.Token, lhs ast.LoxValue, rhs ast.LoxValue) {
	if !isNumber(lhs) || !isNumber(rhs) {
		panic(loxerror.RuntimeError{
			// TODO: print lox types instead of go types
			Message: fmt.Sprintf("Operator '%v': incompatible types %v and %v", operator.Lexeme, reflect.TypeOf(lhs), reflect.TypeOf(rhs)),
		})
	}
}

func assertNumberOperand(operator lexer.Token, rhs ast.LoxValue) {
	if !isNumber(rhs) {
		panic(loxerror.RuntimeError{
			// TODO: print lox types instead of go types
			Message: fmt.Sprintf("Operator '%v': incompatible type %v", operator.Lexeme, reflect.TypeOf(rhs)),
		})
	}
}
