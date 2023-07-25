package ast

import (
	"fmt"

	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
)

// Grammar rules of the Lox language
//
// program -> statement* EOF
//
// declaration -> varDeclaration
//                | funDeclaration
//                | statement
//
// varDeclaration -> IDENTIFIER ( "=" expression )? ";"
//
// funDeclaration -> "fun" function
//
// function -> IDENTIFIER "(" parameters? ")" block
//
// block -> "{" declaration* "}"
// block -> "{" declaration* "}"
// block -> "{" declaration* "}"
//
// statement -> expressionStatement
//              | forStatement
//              | ifStatement
//              | printStatement
//              | whileStatement
//              | block
//
// forStatement -> "for" "(" (varDecl | expression Statement | ";") expression? ";" expression? ")" statement
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
//          | call
//
// call -> primary ( "(" arguments? ")" )*
//
// arguments -> expression ( "," expression )*
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

func (p *Parser) Parse() []Statement {
	statements := make([]Statement, 0)
	for !p.isAtEnd() {
		s := p.declaration()
		statements = append(statements, s)
	}

	return statements
}

func (p *Parser) declaration() Statement {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(loxerror.ParserError); ok {
				p.synchronize()
				return
			}
			panic(r)
		}
	}()

	var statement Statement
	switch {
	case p.match(lexer.Var):
		statement = p.varDeclaration()
	case p.match(lexer.Fun):
		statement = p.function("function")
	case p.match(lexer.LeftBrace):
		statement = NewBlockStatement(p.block())
	default:
		statement = p.statement()
	}

	return statement
}

func (p *Parser) varDeclaration() Statement {
	name := p.consume(lexer.Identifier, "Expect a variable name")
	var initializer Expression
	if p.match(lexer.Equal) {
		initializer = p.expression()
	}

	p.consume(lexer.Semicolon, "Expect ';' after variable declaration")

	return NewVariableStatement(name, initializer)
}

func (p *Parser) function(kind string) Statement {
	name := p.consume(lexer.Identifier, fmt.Sprintf("Expect %s name.", kind))

	p.consume(lexer.LeftParenthesis, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]lexer.Token, 0)
	if !p.check(lexer.RightParenthesis) {
		for next := true; next; next = p.match(lexer.Comma) {
			if len(parameters) >= 255 {
				loxerror.Error(p.peek().Line, "Can't have more than 255 parameters")
			}

			parameters = append(parameters, p.consume(lexer.Identifier, "Expect parameter name"))
		}
	}

	p.consume(lexer.RightParenthesis, "Expect ')' after parameters")
	p.consume(lexer.LeftBrace, "Expect '{' after parameters")
	body := p.block()

	return NewFunctionStatement(name, parameters, body)
}

func (p *Parser) block() []Statement {
	statements := make([]Statement, 0)
	for !p.check(lexer.RightBrace) && !p.isAtEnd() {
		s := p.declaration()
		statements = append(statements, s)
	}
	p.consume(lexer.RightBrace, "Expect '}' after block")

	return statements
}

func (p *Parser) statement() Statement {
	switch {
	case p.match(lexer.For):
		return p.forStatement()
	case p.match(lexer.Print):
		return p.printStatement()
	case p.match(lexer.If):
		return p.ifStatement()
	case p.match(lexer.While):
		return p.whileStatment()
	case p.match(lexer.LeftBrace):
		return NewBlockStatement(p.block())
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) forStatement() Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'for'")
	var initializer Statement
	switch {
	case p.match(lexer.Semicolon):
		// Do nothing
		initializer = nil
	case p.match(lexer.Var):
		initializer = p.varDeclaration()
	default:
		initializer = p.expressionStatement()
	}

	var condition Expression
	if !p.check(lexer.Semicolon) {
		condition = p.expression()
	}
	p.consume(lexer.Semicolon, "Expect ';' after loop condition")

	var increment Expression
	if !p.check(lexer.RightParenthesis) {
		increment = p.expression()
	}
	p.consume(lexer.RightParenthesis, "Expect ')' after for clauses")

	body := p.statement()

	if increment != nil {
		body = NewBlockStatement([]Statement{body, NewExpressionStatement(increment)})
	}

	if condition == nil {
		condition = NewLiteralExpression(NewBooleanValue(true))
	}

	newBody := NewWhileStatement(condition, body)
	var forLoop Statement
	if initializer != nil {
		forLoop = NewBlockStatement([]Statement{initializer, newBody})
	}

	return forLoop
}

func (p *Parser) printStatement() Statement {
	value := p.expression()
	p.consume(lexer.Semicolon, "Expect ';' after value")

	return NewPrintStatement(value)
}

func (p *Parser) ifStatement() Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'if' statement")
	expr := p.expression()
	p.consume(lexer.RightParenthesis, "Expect ')' after 'if' condition")
	thenCode := p.statement()
	var elseCode Statement
	if p.match(lexer.Else) {
		elseCode = p.statement()
	}

	return NewIfStatment(expr, thenCode, elseCode)
}

func (p *Parser) whileStatment() Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'while'")
	condition := p.expression()
	p.consume(lexer.RightParenthesis, "Expect ')' after 'while' condition")
	body := p.statement()

	return NewWhileStatement(condition, body)
}

func (p *Parser) expressionStatement() Statement {
	expression := p.expression()
	p.consume(lexer.Semicolon, "Expect ';' after value")

	return NewExpressionStatement(expression)
}

func (p *Parser) expression() Expression {
	return p.assignment()
}

func (p *Parser) assignment() Expression {
	expression := p.or()

	if p.match(lexer.Equal) {
		equalsToken := p.previous()
		value := p.assignment()

		if expression, ok := expression.(*VariableExpression); ok {
			return NewAssignmentExpression(expression.Name, value)
		}

		// No need to go to put the parser in recovery mode
		// TODO: better error message
		loxerror.Error(equalsToken.Line, "Invalid assignment target")

		return nil
	}

	return expression
}

func (p *Parser) or() Expression {
	expr := p.and()

	for p.match(lexer.Or) {
		operator := p.previous()
		rhs := p.and()
		expr = NewLogicalExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) and() Expression {
	expr := p.equality()

	for p.match(lexer.And) {
		operator := p.previous()
		rhs := p.equality()
		expr = NewLogicalExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) equality() Expression {
	expr := p.comparison()

	for p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.previous()
		rhs := p.comparison()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()

	for p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.previous()
		rhs := p.term()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()

	for p.match(lexer.Plus, lexer.Dash) {
		operator := p.previous()
		rhs := p.factor()

		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()

	for p.match(lexer.Star, lexer.Slash) {
		operator := p.previous()
		rhs := p.unary()
		expr = NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) unary() Expression {
	if p.match(lexer.Bang, lexer.Dash) {
		operator := p.previous()
		rhs := p.unary()

		return NewUnaryExpression(operator, rhs)
	}

	return p.call()
}

func (p *Parser) call() Expression {
	expr := p.primary()

	for {
		if p.match(lexer.LeftParenthesis) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expression) Expression {
	const MaxArgs = 255

	args := make([]Expression, 0)
	if !p.check(lexer.RightParenthesis) {
		for next := true; next; next = p.match(lexer.Comma) {
			if len(args) >= MaxArgs {
				loxerror.Error(p.peek().Line, "Can't have more than 255 arguments")
			}
			args = append(args, p.expression())
		}
	}

	position := p.consume(lexer.RightParenthesis, "Execpect ')' after function arguments")

	return NewCallExpression(callee, position, args)
}

func (p *Parser) primary() Expression {
	// TODO we should probably have BooleanLiteral and maybe ObjectLiteral rather than strings
	switch {
	case p.match(lexer.False):
		return NewLiteralExpression(NewBooleanValue(false))
	case p.match(lexer.True):
		return NewLiteralExpression(NewBooleanValue(true))
	case p.match(lexer.Nil):
		return NewLiteralExpression(&ObjectValue{Value: nil})
	}

	// TODO: is there a better solution?
	if p.match(lexer.Number) {
		return NewLiteralExpression(&NumberValue{Value: p.previous().Literal.(*lexer.NumberLiteral).Value})
	} else if p.match(lexer.String) {
		return NewLiteralExpression(&StringValue{Value: p.previous().Literal.(*lexer.StringLiteral).Value})
	}

	if p.match(lexer.Identifier) {
		return NewVariableExpression(p.previous())
	}

	if p.match(lexer.LeftParenthesis) {
		expr := p.expression()
		p.consume(lexer.RightParenthesis, "Expect ')' after expression")

		return NewGroupingExpression(expr)
	}

	loxerror.Error(p.peek().Line, "expect expression")
	panic(loxerror.ParserError{Message: "expect expression"})
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
func (p *Parser) consume(tokenType lexer.TokenType, message string) lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	loxerror.Error(p.peek().Line, message)
	panic(loxerror.ParserError{Message: message})
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
