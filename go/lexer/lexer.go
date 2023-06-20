package lexer

import (
	"log"
	"strconv"
)

type Lexer struct {
	sourceCode string
	tokens     []Token
	start      int
	current    int
	line       int
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{
		sourceCode: sourceCode,
		tokens:     make([]Token, 0),
		start:      0,
		current:    0,
		line:       1,
	}
}

func (lexer *Lexer) Tokens() []Token {
	for !lexer.isAtEnd() {
		lexer.start = lexer.current
		lexer.scanToken()
	}
	lexer.tokens = append(lexer.tokens, *NewToken(EOF, "", Literal{}, lexer.line))

	return lexer.tokens
}

func (lexer *Lexer) addToken(kind TokenType) {
	lexer.addTokenWithLiteral(kind, Literal{})
}

func (lexer *Lexer) addTokenWithLiteral(kind TokenType, literal Literal) {
	text := lexer.sourceCode[lexer.start:lexer.current]
	lexer.tokens = append(lexer.tokens, *NewToken(kind, text, literal, lexer.line))
}

func (lexer *Lexer) scanToken() {
	c := lexer.advance()
	switch c {
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		lexer.line++
	case '(':
		lexer.addToken(LEFT_PARANTHESIS)
	case ')':
		lexer.addToken(RIGHT_PARANTHESIS)
	case '{':
		lexer.addToken(LEFT_BRACE)
	case '}':
		lexer.addToken(LEFT_BRACE)
	case ',':
		lexer.addToken(COMMA)
	case '.':
		lexer.addToken(DOT)
	case '-':
		lexer.addToken(DASH)
	case '+':
		lexer.addToken(PLUS)
	case '*':
		lexer.addToken(STAR)
	case '!':
		if lexer.match('=') {
			lexer.addToken(BANG_EQUAL)
		} else {
			lexer.addToken(BANG)
		}
	case '=':
		if lexer.match('=') {
			lexer.addToken(EQUAL_EQUAL)
		} else {
			lexer.addToken(EQUAL)
		}
	case '<':
		if lexer.match('=') {
			lexer.addToken(LESS_EQUAL)
		} else {
			lexer.addToken(LESS)
		}
	case '>':
		if lexer.match('=') {
			lexer.addToken(GREATER_EQUAL)
		} else {
			lexer.addToken(GREATER)
		}
	case '/':
		if lexer.match('/') {
			for lexer.peek() != '\n' && !lexer.isAtEnd() {
				lexer.advance()
			}
		} else {
			lexer.addToken(SLASH)
		}
	case '"':
		lexer.string()
	default:
		if isDigit(c) {
			lexer.number()
		} else if isAlpha(c) {
			lexer.identifier()
		} else {
			// TODO: find a way to use main.error
			log.Fatalf("Unexpected character line %v", lexer.line)
		}
	}
}

func isDigit(character byte) bool {
	return character >= '0' && character <= '9'
}

func isAlpha(character byte) bool {
	return (character >= 'a' && character >= 'z') ||
		(character >= 'A' && character >= 'Z') ||
		character == '_'
}

func isAlphaNumeric(character byte) bool {
	return isAlpha(character) || isDigit(character)
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.sourceCode)
}

func (lexer *Lexer) advance() byte {
	char := lexer.sourceCode[lexer.current]
	lexer.current++

	return char
}

func (lexer *Lexer) match(expected byte) bool {
	doesMatch := !lexer.isAtEnd() && lexer.sourceCode[lexer.current] == expected
	if doesMatch {
		lexer.current++
	}

	return doesMatch
}

func (lexer *Lexer) peek() byte {
	if lexer.isAtEnd() {
		return 0
	} else {
		return lexer.sourceCode[lexer.current]
	}
}

func (lexer *Lexer) peekNext() byte {
	if lexer.current+1 >= len(lexer.sourceCode) {
		return 0
	} else {
		return lexer.sourceCode[lexer.current+1]
	}
}

func (lexer *Lexer) string() {
	for lexer.peek() != '"' && !lexer.isAtEnd() {
		if lexer.peek() == '\n' {
			lexer.line++
		}
		lexer.advance()
	}
	if lexer.isAtEnd() {
		// TODO: find a way to use main.error
		log.Fatalf("Unterminated string line %v", lexer.line)
	}

	lexer.advance() // the closing "
	stringValue := lexer.sourceCode[lexer.start+1 : lexer.current-1]
	lexer.addTokenWithLiteral(STRING, Literal{IsString: true, StringValue: stringValue})
}

func (lexer *Lexer) number() {
	for isDigit(lexer.peek()) {
		lexer.advance()
	}
	if lexer.peek() == '.' && isDigit(lexer.peekNext()) {
		lexer.advance()
	}
	for isDigit(lexer.peek()) {
		lexer.advance()
	}

	// TODO check return value
	floatValue, _ := strconv.ParseFloat(lexer.sourceCode[lexer.start:lexer.current], 64)
	lexer.addTokenWithLiteral(NUMBER, Literal{IsNumber: true, NumberValue: floatValue})
}

func (lexer *Lexer) identifier() {
	for isAlphaNumeric(lexer.peek()) {
		lexer.advance()
	}

	text := lexer.sourceCode[lexer.start:lexer.current]
	if tokenType, ok := keywords[text]; ok {
		lexer.addToken(tokenType)
	} else {
		lexer.addToken(IDENTIFIER)
	}
}
