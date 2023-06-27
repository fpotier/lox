package ast

import (
	"errors"
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

// program -> statement* EOF
func (parser *Parser) Parse() ([]Statement, error) {
	statements := make([]Statement, 0)
	for !parser.isAtEnd() {
		// TODO: handle error
		s, err := parser.declaration()
		if err != nil {
			fmt.Println(err)
		} else {
			statements = append(statements, s)
		}
	}

	return statements, nil
}

// declaration -> varDeclaration
//
//	| statement
//	| block
func (parser *Parser) declaration() (Statement, error) {
	var statement Statement
	var err error
	if parser.match(lexer.VAR) {
		statement, err = parser.varDeclaration()
	} else if parser.match(lexer.LEFT_BRACE) {
		statements, innerErr := parser.block()
		if innerErr == nil {
			statement = NewBlockStatement(statements)
		}
		err = innerErr
	} else {
		statement, err = parser.statement()
	}
	if err != nil {
		parser.synchronize()
		return nil, err
	}

	return statement, nil
}

// varDeclaration -> IDENTIFIER ( "=" expression )? ";"
func (parser *Parser) varDeclaration() (Statement, error) {
	// TODO: better error message
	name, err := parser.consume(lexer.IDENTIFIER, "Expect a variable name")
	if err != nil {
		return nil, err
	}
	var initializer Expression
	if parser.match(lexer.EQUAL) {
		initializer, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	// TODO: better error message
	_, err = parser.consume(lexer.SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return NewVariableStatement(name, initializer), nil
}

// block -> "{" declaration* "}"
func (parser *Parser) block() ([]Statement, error) {
	statements := make([]Statement, 0)
	for !parser.check(lexer.RIGHT_BRACE) && !parser.isAtEnd() {
		s, err := parser.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}
	_, err := parser.consume(lexer.RIGHT_BRACE, "Expect '}' after block")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

// statement -> expressionStatement
//
//	| printStatement
func (parser *Parser) statement() (Statement, error) {
	if parser.match(lexer.PRINT) {
		return parser.printStatement()
	} else {
		return parser.expressionStatement()
	}
}

// printStatement -> "print" expression ";"
func (parser *Parser) printStatement() (Statement, error) {
	value, err := parser.expression()
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(lexer.SEMICOLON, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return NewPrintStatement(value), nil
}

// expressionStatement -> expression ";"
func (parser *Parser) expressionStatement() (Statement, error) {
	expression, err := parser.expression()
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(lexer.SEMICOLON, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return NewExpressionStatement(expression), nil
}

// expression -> assignment
func (parser *Parser) expression() (Expression, error) {
	return parser.assignment()
}

// assignment -> IDENTIFIER "=" assignment
//
//	| equality
func (parser *Parser) assignment() (Expression, error) {
	expression, err := parser.equality()
	if err != nil {
		return nil, err
	}

	if parser.match(lexer.EQUAL) {
		equalsToken := parser.previous()
		value, err := parser.assignment()
		if err != nil {
			return nil, err
		}

		if expression, ok := expression.(*VariableExpression); ok {
			return NewAssignmentExpression(expression.Name, value), nil
		}

		// No need to go to put the parser in recovery mode
		// TODO: better error message
		loxerror.Error(equalsToken.Line, "Invalid assignment target")

		return nil, errors.New("invalid assignment target")
	}

	return expression, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison ) *
func (parser *Parser) equality() (Expression, error) {
	expr, err := parser.comparison()
	if err != nil {
		return nil, err
	}

	for parser.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := parser.previous()
		rhs, err := parser.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (parser *Parser) comparison() (Expression, error) {
	expr, err := parser.term()
	if err != nil {
		return nil, err
	}

	for parser.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := parser.previous()
		rhs, err := parser.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

// term -> factor ( ( "-" | "+" ) comparison )*
func (parser *Parser) term() (Expression, error) {
	expr, err := parser.factor()
	if err != nil {
		return nil, err
	}

	for parser.match(lexer.PLUS, lexer.DASH) {
		operator := parser.previous()
		rhs, err := parser.factor()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

// factor -> unary ( ( "/" | "*" ) unary )*
func (parser *Parser) factor() (Expression, error) {
	expr, err := parser.unary()
	if err != nil {
		return nil, err
	}

	for parser.match(lexer.STAR, lexer.SLASH) {
		operator := parser.previous()
		rhs, err := parser.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary
// | primary
func (parser *Parser) unary() (Expression, error) {
	if parser.match(lexer.BANG, lexer.DASH) {
		operator := parser.previous()
		rhs, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpression(operator, rhs), nil
	}

	return parser.primary()
}

// primary -> NUMBER
// | STRING
// | "true"
// | "false"
// | "nil"
// | "(" expression ")"
// | IDENTIFIER
func (parser *Parser) primary() (Expression, error) {
	// TODO we should probably have BooleanLiteral and maybe ObjectLiteral rather than strings
	if parser.match(lexer.FALSE) {
		return NewLiteralExpression(NewBooleanValue(false)), nil
	} else if parser.match(lexer.TRUE) {
		return NewLiteralExpression(NewBooleanValue(true)), nil
	} else if parser.match(lexer.NIL) {
		return NewLiteralExpression(&ObjectValue{Value: nil}), nil
	}

	// TODO: is there a better solution?
	if parser.match(lexer.NUMBER) {
		return NewLiteralExpression(&NumberValue{Value: parser.previous().Literal.(*lexer.NumberLiteral).Value}), nil
	} else if parser.match(lexer.STRING) {
		return NewLiteralExpression(&StringValue{Value: parser.previous().Literal.(*lexer.StringLiteral).Value}), nil
	}

	if parser.match(lexer.IDENTIFIER) {
		return NewVariableExpression(parser.previous()), nil
	}

	if parser.match(lexer.LEFT_PARANTHESIS) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		_, err = parser.consume(lexer.RIGHT_PARANTHESIS, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return NewGroupingExpression(expr), nil
	}

	return nil, errors.New("expect expression")
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) check(tokenType lexer.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tokenType
}

func (parser *Parser) advance() lexer.Token {
	if !parser.isAtEnd() {
		parser.current++
	}
	return parser.previous()
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == lexer.EOF
}

func (parser *Parser) peek() lexer.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) previous() lexer.Token {
	return parser.tokens[parser.current-1]
}

// TODO: all call site should have better error messages
func (parser *Parser) consume(tokenType lexer.TokenType, message string) (lexer.Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	} else {
		return parser.peek(), errors.New(message)
	}
}

func (parser *Parser) synchronize() {
	parser.advance()

	for !parser.isAtEnd() {
		switch parser.previous().Type {
		case lexer.SEMICOLON:
		case lexer.CLASS:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}

		parser.advance()
	}
}
