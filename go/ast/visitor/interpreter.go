package visitor

import (
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Interpreter struct {
	// value can be virtually anything (string, number, boolean, object, nil, etc.)
	value ast.LoxValue
}

func (visitor *Interpreter) Eval(expression ast.Expression) interface{} {
	// This allows us to emulate a try/catch mechanism to exit the visitor as soon as possible
	// without changing the Visit...() methods to return an error and propagate manually these errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(loxerror.RuntimeError); ok {
				fmt.Println(err)
				return
			}
			panic(r)
		}
	}()
	expression.Accept(visitor)
	return visitor.value
}

func (visitor *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := visitor.evaluate(binaryExpression.Lhs)
	rhs := visitor.evaluate(binaryExpression.Rhs)

	switch binaryExpression.Operator.Type {
	case lexer.PLUS:
		switch lhs := lhs.(type) {
		case *ast.NumberValue:
			if rhs, ok := rhs.(*ast.NumberValue); ok {
				visitor.value = &ast.NumberValue{Value: lhs.Value + rhs.Value}
			} else {
				panic(loxerror.RuntimeError{Message: "invalid operand type: Number + String"})
			}
		case *ast.StringValue:
			if rhs, ok := rhs.(*ast.StringValue); ok {
				visitor.value = &ast.StringValue{Value: lhs.Value + rhs.Value}
			}
		}
	case lexer.DASH:
		visitor.value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value - rhs.(*ast.NumberValue).Value}
	case lexer.STAR:
		visitor.value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value * rhs.(*ast.NumberValue).Value}
	case lexer.SLASH:
		visitor.value = &ast.NumberValue{Value: lhs.(*ast.NumberValue).Value / rhs.(*ast.NumberValue).Value}
	}
}

func (visitor *Interpreter) VisitGroupingExpression(groupingExpression *ast.GroupingExpression) {
	visitor.value = visitor.evaluate(groupingExpression.Expr)
}

func (visitor *Interpreter) VisitLiteralExpression(literalExpression *ast.LiteralExpression) {
	visitor.value = literalExpression.LoxValue()
}

func (visitor *Interpreter) VisitUnaryExpression(unaryExpression *ast.UnaryExpression) {
	rhs := visitor.evaluate(unaryExpression.Rhs)

	switch unaryExpression.Operator.Type {
	case lexer.BANG:
		visitor.value = &ast.BooleanValue{Value: !rhs.IsTruthy()}
	case lexer.DASH:
		visitor.value = &ast.NumberValue{Value: -rhs.(*ast.NumberValue).Value}
	}
}

func (visitor *Interpreter) evaluate(expression ast.Expression) ast.LoxValue {
	newVisitor := &Interpreter{}
	expression.Accept(newVisitor)

	return newVisitor.value
}
