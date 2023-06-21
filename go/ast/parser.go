package ast

import (
	"errors"

	"github.com/fpotier/crafting-interpreters/go/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func (parser *Parser) Parse() Expression {
	// TODO error handling
	return parser.expression()
}

// expression -> equality
func (parser *Parser) expression() Expression {
	return parser.equality()
}

//equality -> comparison ( ( "!=" | "==" ) comparison ) *
func (parser *Parser) equality() Expression {
	expr := parser.comparison()

	for parser.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := parser.previous()
		rhs := parser.comparison()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (parser *Parser) comparison() Expression {
	expr := parser.term()

	for parser.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := parser.previous()
		rhs := parser.term()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

// term -> factor ( ( "-" | "+" ) comparison )*
func (parser *Parser) term() Expression {
	expr := parser.factor()

	for parser.match(lexer.PLUS, lexer.DASH) {
		operator := parser.previous()
		rhs := parser.factor()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

// factor -> unary ( ( "/" | "*" ) unary )*
func (parser *Parser) factor() Expression {
	expr := parser.unary()

	for parser.match(lexer.STAR, lexer.SLASH) {
		operator := parser.previous()
		rhs := parser.unary()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

// unary -> ( "!" | "-" ) unary | primary
func (parser *Parser) unary() Expression {
	if parser.match(lexer.BANG, lexer.DASH) {
		operator := parser.previous()
		rhs := parser.unary()
		return NewUnaryExpression(operator, rhs)
	}

	return parser.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (parser *Parser) primary() Expression {
	// TODO we should probably have BooleanLiteral and maybe ObjectLiteral rather than strings
	if parser.match(lexer.FALSE) {
		return NewLiteralExpression(&lexer.StringLiteral{Value: "false"})
	} else if parser.match(lexer.TRUE) {
		return NewLiteralExpression(&lexer.StringLiteral{Value: "true"})
	} else if parser.match(lexer.NIL) {
		return NewLiteralExpression(&lexer.StringLiteral{Value: "nil"})
	}

	if parser.match(lexer.NUMBER, lexer.STRING) {
		return NewLiteralExpression(parser.previous().Literal)
	}

	if parser.match(lexer.LEFT_PARANTHESIS) {
		expr := parser.expression()
		// TODO handle error
		parser.consume(lexer.RIGHT_PARANTHESIS, "Expect ')' after expression")
		return NewGroupingExpression(expr)
	}

	panic("Never reached")
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

func (parser *Parser) consume(tokenType lexer.TokenType, message string) (lexer.Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	} else {
		return parser.peek(), errors.New(message)
	}
}
