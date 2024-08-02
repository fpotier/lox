package parser

import (
	"fmt"

	"github.com/fpotier/lox/go/pkg/ast"
	"github.com/fpotier/lox/go/pkg/lexer"
	"github.com/fpotier/lox/go/pkg/loxerror"
)

// Grammar rules of the Lox language
//
// program -> declaration* EOF
//
// declaration -> varDeclaration
//                | funDeclaration
//                | classDeclaration
//                | statement
//
// varDeclaration -> IDENTIFIER ( "=" expression )? ";"
//
// funDeclaration -> "fun" function
//
// classDeclaration -> "class" IDENTIFIER ( "<"  IDENTIFIER )? "{" function* "}"
//
// statement -> expressionStatement
//              | forStatement
//              | ifStatement
//              | printStatement
//              | returnStatement
//              | whileStatement
//              | block
//
// expressionStatement -> expression ";"
//
// forStatement -> "for" "(" (varDecl | expression Statement | ";") expression? ";" expression? ")" statement
//
// ifStatement -> "if" "(" expression ")" ( "else" statement )?
//
// printStatement -> "print" expression ";"
//
// returnStatement -> "return" expression? ";"
//
// whileStatment -> "while" "(" expression ")" statement
//
// block -> "{" declaration* "}"
//
// expression -> assignment
//
// assignment -> ( call "." )? IDENTIFIER "=" assignment
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
// call -> primary ( "(" arguments? ")" | "." IDENTIFIER )*
//
// primary -> NUMBER
//            | STRING
//            | "true"
//            | "false"
//            | "nil"
//            | "(" expression ")"
//            | IDENTIFIER
//            | "super" "." IDENTIFIER
//
// function -> IDENTIFIER "(" parameters? ")" block
//
// parameters -> expression ( "," IDENTIFIER )*
//
// arguments -> expression ( "," expression )*
//

func NewParser(errorFormatter loxerror.ErrorFormatter, tokens []lexer.Token) *Parser {
	return &Parser{
		errorFormatter: errorFormatter,
		tokens:         tokens,
		current:        0,
	}
}

func (p *Parser) Parse() []ast.Statement {
	statements := make([]ast.Statement, 0)
	for !p.isAtEnd() {
		s := p.declaration()
		statements = append(statements, s)
	}

	return statements
}

func (p *Parser) declaration() ast.Statement {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*ParseError); ok {
				p.synchronize()
				return
			}
			panic(r)
		}
	}()

	var statement ast.Statement
	switch {
	case p.match(lexer.Var):
		statement = p.varDeclaration()
	case p.match(lexer.Fun):
		statement = p.function("function")
	case p.match(lexer.Class):
		statement = p.classDeclaration()
	case p.match(lexer.LeftBrace):
		statement = ast.NewBlockStatement(p.block())
	default:
		statement = p.statement()
	}

	return statement
}

func (p *Parser) varDeclaration() ast.Statement {
	name := p.consume(lexer.Identifier, "Expect a variable name")
	var initializer ast.Expression
	if p.match(lexer.Equal) {
		initializer = p.expression()
	}

	p.consume(lexer.Semicolon, "Expect ';' after variable declaration")

	return ast.NewVariableStatement(name, initializer)
}

func (p *Parser) function(kind string) ast.Statement {
	name := p.consume(lexer.Identifier, fmt.Sprintf("Expect %s name.", kind))

	p.consume(lexer.LeftParenthesis, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]lexer.Token, 0)
	if !p.check(lexer.RightParenthesis) {
		for next := true; next; next = p.match(lexer.Comma) {
			if len(parameters) >= ast.Limits.MaxArgs {
				p.errorFormatter.PushError(NewParseError(p.peek().Line, "Can't have more than 255 parameters"))
			}

			parameters = append(parameters, p.consume(lexer.Identifier, "Expect parameter name"))
		}
	}

	p.consume(lexer.RightParenthesis, "Expect ')' after parameters")
	p.consume(lexer.LeftBrace, "Function body must start with '{'")
	body := p.block()

	return ast.NewFunctionStatement(name, parameters, body)
}

func (p *Parser) classDeclaration() ast.Statement {
	name := p.consume(lexer.Identifier, "Expect class name")

	var superclass *ast.VariableExpression
	if p.match(lexer.Less) {
		p.consume(lexer.Identifier, "Expect superclass name")
		superclass = ast.NewVariableExpression(p.previous())
	}

	p.consume(lexer.LeftBrace, "Expect '{' before class body")

	methods := make([]*ast.FunctionStatement, 0)
	for !p.check(lexer.RightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method").(*ast.FunctionStatement))
	}

	p.consume(lexer.RightBrace, "Expect '}' after class body")

	return ast.NewClassStatement(name, superclass, methods)
}

func (p *Parser) statement() ast.Statement {
	switch {
	case p.match(lexer.For):
		return p.forStatement()
	case p.match(lexer.Print):
		return p.printStatement()
	case p.match(lexer.If):
		return p.ifStatement()
	case p.match(lexer.Return):
		return p.returnStatement()
	case p.match(lexer.While):
		return p.whileStatment()
	case p.match(lexer.LeftBrace):
		return ast.NewBlockStatement(p.block())
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) expressionStatement() ast.Statement {
	expression := p.expression()
	p.consume(lexer.Semicolon, "Expect ';' after value")

	return ast.NewExpressionStatement(expression)
}

// TODO: check where we create an additional block
func (p *Parser) forStatement() ast.Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'for'")
	var initializer ast.Statement
	switch {
	case p.match(lexer.Semicolon):
		// Do nothing
		initializer = nil
	case p.match(lexer.Var):
		initializer = p.varDeclaration()
	default:
		initializer = p.expressionStatement()
	}

	var condition ast.Expression
	if !p.check(lexer.Semicolon) {
		condition = p.expression()
	}
	p.consume(lexer.Semicolon, "Expect ';' after loop condition")

	var increment ast.Expression
	if !p.check(lexer.RightParenthesis) {
		increment = p.expression()
	}
	p.consume(lexer.RightParenthesis, "Expect ')' after for clauses")

	body := p.statement()

	if increment != nil {
		// Note: I think this creates a useless block
		body = ast.NewBlockStatement([]ast.Statement{body, ast.NewExpressionStatement(increment)})
	}

	if condition == nil {
		condition = ast.NewLiteralExpression(ast.NewBooleanValue(true))
	}

	var forLoop ast.Statement = ast.NewWhileStatement(condition, body)
	if initializer != nil {
		forLoop = ast.NewBlockStatement([]ast.Statement{initializer, forLoop})
	}

	return forLoop
}

func (p *Parser) ifStatement() ast.Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'if' statement")
	expr := p.expression()
	p.consume(lexer.RightParenthesis, "Expect ')' after 'if' condition")
	thenCode := p.statement()
	var elseCode ast.Statement
	if p.match(lexer.Else) {
		elseCode = p.statement()
	}

	return ast.NewIfStatment(expr, thenCode, elseCode)
}

func (p *Parser) printStatement() ast.Statement {
	value := p.expression()
	p.consume(lexer.Semicolon, "Expect ';' after value")

	return ast.NewPrintStatement(value)
}

func (p *Parser) returnStatement() ast.Statement {
	keyword := p.previous()
	var value ast.Expression

	if !p.check(lexer.Semicolon) {
		value = p.expression()
	}

	p.consume(lexer.Semicolon, "Expect ';' after return value")

	return ast.NewReturnStatement(keyword, value)
}

func (p *Parser) whileStatment() ast.Statement {
	p.consume(lexer.LeftParenthesis, "Expect '(' after 'while'")
	condition := p.expression()
	p.consume(lexer.RightParenthesis, "Expect ')' after 'while' condition")
	body := p.statement()

	return ast.NewWhileStatement(condition, body)
}

func (p *Parser) block() []ast.Statement {
	statements := make([]ast.Statement, 0)
	for !p.check(lexer.RightBrace) && !p.isAtEnd() {
		s := p.declaration()
		statements = append(statements, s)
	}
	p.consume(lexer.RightBrace, "Expect '}' after block")

	return statements
}

func (p *Parser) expression() ast.Expression {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expression {
	expression := p.or()

	if p.match(lexer.Equal) {
		equalsToken := p.previous()
		value := p.assignment()

		switch expression := expression.(type) {
		case *ast.VariableExpression:
			return ast.NewAssignmentExpression(expression.Name, value)
		case *ast.GetExpression:
			return ast.NewSetExpression(expression.Object, expression.Name, value)
		}

		// No need to go to put the parser in recovery mode
		// TODO: better error message
		p.errorFormatter.PushError(NewParseError(equalsToken.Line, "Invalid assignment target"))

		return nil
	}

	return expression
}

func (p *Parser) or() ast.Expression {
	expr := p.and()

	for p.match(lexer.Or) {
		operator := p.previous()
		rhs := p.and()
		expr = ast.NewLogicalExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) and() ast.Expression {
	expr := p.equality()

	for p.match(lexer.And) {
		operator := p.previous()
		rhs := p.equality()
		expr = ast.NewLogicalExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) equality() ast.Expression {
	expr := p.comparison()

	for p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.previous()
		rhs := p.comparison()
		expr = ast.NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) comparison() ast.Expression {
	expr := p.term()

	for p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.previous()
		rhs := p.term()
		expr = ast.NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) term() ast.Expression {
	expr := p.factor()

	for p.match(lexer.Plus, lexer.Dash) {
		operator := p.previous()
		rhs := p.factor()

		expr = ast.NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) factor() ast.Expression {
	expr := p.unary()

	for p.match(lexer.Star, lexer.Slash) {
		operator := p.previous()
		rhs := p.unary()
		expr = ast.NewBinaryExpression(expr, operator, rhs)
	}

	return expr
}

func (p *Parser) unary() ast.Expression {
	if p.match(lexer.Bang, lexer.Dash) {
		operator := p.previous()
		rhs := p.unary()

		return ast.NewUnaryExpression(operator, rhs)
	}

	return p.call()
}

func (p *Parser) call() ast.Expression {
	expr := p.primary()

loop:
	for {
		switch {
		case p.match(lexer.LeftParenthesis):
			expr = p.finishCall(expr)
		case p.match(lexer.Dot):
			name := p.consume(lexer.Identifier, "Expect property name after '.'")
			expr = ast.NewGetExpression(expr, name)
		default:
			break loop
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expression) ast.Expression {
	args := make([]ast.Expression, 0)
	if !p.check(lexer.RightParenthesis) {
		for next := true; next; next = p.match(lexer.Comma) {
			if len(args) >= ast.Limits.MaxArgs {
				p.errorFormatter.PushError(NewParseError(p.peek().Line, "Can't have more than 255 arguments"))
			}
			args = append(args, p.expression())
		}
	}

	position := p.consume(lexer.RightParenthesis, "Expect ')' after function arguments")

	return ast.NewCallExpression(callee, position, args)
}

func (p *Parser) primary() ast.Expression {
	switch {
	case p.match(lexer.False):
		return ast.NewLiteralExpression(ast.NewBooleanValue(false))
	case p.match(lexer.True):
		return ast.NewLiteralExpression(ast.NewBooleanValue(true))
	case p.match(lexer.Nil):
		return ast.NewLiteralExpression(ast.NewNilValue())
	case p.match(lexer.This):
		return ast.NewThisExpression(p.previous())
	case p.match(lexer.Super):
		keyword := p.previous()
		p.consume(lexer.Dot, "Expect '.' after 'super'")
		method := p.consume(lexer.Identifier, "Expect superclass method name")
		return ast.NewSuperExpression(keyword, method)
	case p.match(lexer.Number):
		return ast.NewLiteralExpression(ast.NewNumberValue(p.previous().Literal.(*lexer.NumberLiteral).Value))
	case p.match(lexer.String):
		return ast.NewLiteralExpression(ast.NewStringValue(p.previous().Literal.(*lexer.StringLiteral).Value))
	case p.match(lexer.Identifier):
		return ast.NewVariableExpression(p.previous())
	case p.match(lexer.LeftParenthesis):
		expr := p.expression()
		p.consume(lexer.RightParenthesis, "Expect ')' after expression")
		return ast.NewGroupingExpression(expr)
	default:
		err := NewParseError(p.peek().Line, "expect expression")
		p.errorFormatter.PushError(err)
		panic(err)
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

	err := NewParseError(p.peek().Line, message)
	p.errorFormatter.PushError(err)
	panic(err)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == lexer.Semicolon {
			return
		}

		switch p.peek().Type {
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

type Parser struct {
	errorFormatter loxerror.ErrorFormatter
	tokens         []lexer.Token
	current        int
}
