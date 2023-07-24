package ast

import (
	"errors"
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

// Grammar rules of the Lox language
//
// program -> statement* EOF
//
// declaration -> varDeclaration
//                | statement
//                | block
//
// varDeclaration -> IDENTIFIER ( "=" expression )? ";"
//
// block -> "{" declaration* "}"
//
// statement -> expressionStatement
//              | ifStatement
//              | printStatement
//              | whileStatement
//              | block
//
// ifStatement -> "if" "(" expression ")" ( "else" statement )?
//
// printStatement -> "print" expression ";"
//
// whileStatment -> "while" "(" expression ")" statement
//
// expressionStatement -> expression ";"
//
// expression -> assignment
//
// assignment -> IDENTIFIER "=" assignment
//               | logic_or
//
// logic_or -> logic_and ( "or" logic_and )*
//
// logic_and -> equality ( "and" equality )*
//
// equality -> comparison ( ( "!=" | "==" ) comparison ) *
//
// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
//
// term -> factor ( ( "-" | "+" ) comparison )*
//
// factor -> unary ( ( "/" | "*" ) unary )*
//
// unary -> ( "!" | "-" ) unary
//          | primary
//
// primary -> NUMBER
//            | STRING
//            | "true"
//            | "false"
//            | "nil"
//            | "(" expression ")"
//            | IDENTIFIER

type Parser struct {
	tokens  []lexer.Token
	current int
}

func (p *Parser) Parse() ([]Statement, error) {
	statements := make([]Statement, 0)
	for !p.isAtEnd() {
		// TODO: handle error
		s, err := p.declaration()
		if err != nil {
			fmt.Println(err)
		} else {
			statements = append(statements, s)
		}
	}

	return statements, nil
}

func (p *Parser) declaration() (Statement, error) {
	var statement Statement
	var err error
	switch {
	case p.match(lexer.Var):
		statement, err = p.varDeclaration()
	case p.match(lexer.LeftBrace):
		statements, innerErr := p.block()
		if innerErr == nil {
			statement = NewBlockStatement(statements)
		}
		err = innerErr
	default:
		statement, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, err
	}

	return statement, nil
}

func (p *Parser) varDeclaration() (Statement, error) {
	// TODO: better error message
	name, err := p.consume(lexer.Identifier, "Expect a variable name")
	if err != nil {
		return nil, err
	}
	var initializer Expression
	if p.match(lexer.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	// TODO: better error message
	_, err = p.consume(lexer.Semicolon, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return NewVariableStatement(name, initializer), nil
}

func (p *Parser) block() ([]Statement, error) {
	statements := make([]Statement, 0)
	for !p.check(lexer.RightBrace) && !p.isAtEnd() {
		s, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}
	_, err := p.consume(lexer.RightBrace, "Expect '}' after block")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) statement() (Statement, error) {
	switch {
	case p.match(lexer.Print):
		return p.printStatement()
	case p.match(lexer.If):
		return p.ifStatement()
	case p.match(lexer.While):
		return p.whileStatment()
	case p.match(lexer.LeftBrace):
		statements, err := p.block()
		return NewBlockStatement(statements), err
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(lexer.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return NewPrintStatement(value), nil
}

func (p *Parser) ifStatement() (Statement, error) {
	// TODO error handling
	_, _ = p.consume(lexer.LeftParenthesis, "Expect '(' after 'if' statement")
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	// TODO error handling
	_, _ = p.consume(lexer.RightParenthesis, "Expect ')' after 'if' condition")
	thenCode, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseCode Statement
	if p.match(lexer.Else) {
		elseCode, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIfStatment(expr, thenCode, elseCode), nil
}

func (p *Parser) whileStatment() (Statement, error) {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'while'")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(lexer.RightParenthesis, "Expect ')' after 'while' condition")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return NewWhileStatement(condition, body), nil
}

func (p *Parser) expressionStatement() (Statement, error) {
	expression, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(lexer.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return NewExpressionStatement(expression), nil
}

func (p *Parser) expression() (Expression, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expression, error) {
	expression, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(lexer.Equal) {
		equalsToken := p.previous()
		value, err := p.assignment()
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

func (p *Parser) or() (Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.Or) {
		operator := p.previous()
		rhs, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) and() (Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.And) {
		operator := p.previous()
		rhs, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) equality() (Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.previous()
		rhs, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.previous()
		rhs, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) term() (Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.Plus, lexer.Dash) {
		operator := p.previous()
		rhs, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) factor() (Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.Star, lexer.Slash) {
		operator := p.previous()
		rhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr, nil
}

func (p *Parser) unary() (Expression, error) {
	if p.match(lexer.Bang, lexer.Dash) {
		operator := p.previous()
		rhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpression(operator, rhs), nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expression, error) {
	// TODO we should probably have BooleanLiteral and maybe ObjectLiteral rather than strings
	switch {
	case p.match(lexer.False):
		return NewLiteralExpression(NewBooleanValue(false)), nil
	case p.match(lexer.True):
		return NewLiteralExpression(NewBooleanValue(true)), nil
	case p.match(lexer.Nil):
		return NewLiteralExpression(&ObjectValue{Value: nil}), nil
	}

	// TODO: is there a better solution?
	if p.match(lexer.Number) {
		return NewLiteralExpression(&NumberValue{Value: p.previous().Literal.(*lexer.NumberLiteral).Value}), nil
	} else if p.match(lexer.String) {
		return NewLiteralExpression(&StringValue{Value: p.previous().Literal.(*lexer.StringLiteral).Value}), nil
	}

	if p.match(lexer.Identifier) {
		return NewVariableExpression(p.previous()), nil
	}

	if p.match(lexer.LeftParenthesis) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(lexer.RightParenthesis, "Expect ')' after expression")
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

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.EOF
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

// TODO: all call site should have better error messages
func (p *Parser) consume(tokenType lexer.TokenType, message string) (lexer.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return p.peek(), errors.New(message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		switch p.previous().Type {
		case lexer.Semicolon:
			fallthrough
		case lexer.Class:
			fallthrough
		case lexer.Var:
			fallthrough
		case lexer.For:
			fallthrough
		case lexer.If:
			fallthrough
		case lexer.While:
			fallthrough
		case lexer.Print:
			fallthrough
		case lexer.Return:
			return
		default:
			p.advance()
		}
	}
}
