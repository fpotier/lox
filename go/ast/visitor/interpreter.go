package visitor

import (
	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/lexer"
)

type Interpreter struct {
	// value can be virtually anything (string, number, boolean, object, nil, etc.)
	value ast.LoxValue
}

func (visitor *Interpreter) Eval(expression ast.Expression) interface{} {
	expression.Accept(visitor)
	return visitor.value
}

func (visitor *Interpreter) VisitBinaryExpression(binaryExpression *ast.BinaryExpression) {
	lhs := visitor.evaluate(binaryExpression.Lhs)
	rhs := visitor.evaluate(binaryExpression.Rhs)

	switch binaryExpression.Operator.Type {
	case lexer.DASH:
		visitor.value = lhs.(*ast.NumberValue).Substract(rhs.(*ast.NumberValue))
	case lexer.STAR:
		visitor.value = lhs.(*ast.NumberValue).Multiply(rhs.(*ast.NumberValue))
	case lexer.SLASH:
		visitor.value = lhs.(*ast.NumberValue).Divide(rhs.(*ast.NumberValue))
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

func isTruthy(value interface{}) bool {
	switch value := value.(type) {
	case bool:
		return value
	default:
		return value != nil
	}
}
