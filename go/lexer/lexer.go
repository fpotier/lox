package lexer

import (
	"fmt"
	"strconv"

	"github.com/fpotier/crafting-interpreters/go/loxerror"
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
	lexer.tokens = append(lexer.tokens, *NewToken(EOF, "", nil, lexer.line))

	return lexer.tokens
}

func (lexer *Lexer) addToken(kind TokenType) {
	lexer.addTokenWithLiteral(kind, nil)
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
		lexer.addToken(LeftParenthesis)
	case ')':
		lexer.addToken(RightParenthesis)
	case '{':
		lexer.addToken(LeftBrace)
	case '}':
		lexer.addToken(RightBrace)
	case ';':
		lexer.addToken(Semicolon)
	case ',':
		lexer.addToken(Comma)
	case '.':
		lexer.addToken(Dot)
	case '-':
		lexer.addToken(Dash)
	case '+':
		lexer.addToken(Plus)
	case '*':
		lexer.addToken(Star)
	case '!':
		if lexer.match('=') {
			lexer.addToken(BangEqual)
		} else {
			lexer.addToken(Bang)
		}
	case '=':
		if lexer.match('=') {
			lexer.addToken(EqualEqual)
		} else {
			lexer.addToken(Equal)
		}
	case '<':
		if lexer.match('=') {
			lexer.addToken(LessEqual)
		} else {
			lexer.addToken(Less)
		}
	case '>':
		if lexer.match('=') {
			lexer.addToken(GreaterEqual)
		} else {
			lexer.addToken(Greater)
		}
	case '/':
		if lexer.match('/') {
			for lexer.peek() != '\n' && !lexer.isAtEnd() {
				lexer.advance()
			}
		} else {
			lexer.addToken(Slash)
		}
	case '"':
		lexer.string()
	default:
		switch {
		case isDigit(c):
			lexer.number()
		case isAlpha(c):
			lexer.identifier()
		default:
			loxerror.Error(lexer.line, fmt.Sprintf("Unexpected character '%c'", c))
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
	}

	return lexer.sourceCode[lexer.current]
}

func (lexer *Lexer) peekNext() byte {
	if lexer.current+1 >= len(lexer.sourceCode) {
		return 0
	}

	return lexer.sourceCode[lexer.current+1]
}

func (lexer *Lexer) string() {
	for lexer.peek() != '"' && !lexer.isAtEnd() {
		if lexer.peek() == '\n' {
			lexer.line++
		}

		lexer.advance()
	}

	if lexer.isAtEnd() {
		loxerror.Error(lexer.line, "Unterminated string")
		return
	}

	lexer.advance() // the closing "
	stringValue := lexer.sourceCode[lexer.start+1 : lexer.current-1]
	lexer.addTokenWithLiteral(String, &StringLiteral{Value: stringValue})
}

// Lox number should only be of this form: (\d*)(\.?)(\d*)
func (lexer *Lexer) number() {
	for isDigit(lexer.peek()) {
		lexer.advance()
	}

	if lexer.peek() == '.' && isDigit(lexer.peekNext()) {
		// Consume the '.'
		lexer.advance()

		for isDigit(lexer.peek()) {
			lexer.advance()
		}
	}

	floatValue, err := strconv.ParseFloat(lexer.sourceCode[lexer.start:lexer.current], 64)
	if err != nil {
		loxerror.Error(lexer.line,
			fmt.Sprintf("Error converting %v to float: %v", lexer.sourceCode[lexer.start:lexer.current], err))
		return
	}
	lexer.addTokenWithLiteral(Number, &NumberLiteral{Value: floatValue})
}

func (lexer *Lexer) identifier() {
	for isAlphaNumeric(lexer.peek()) {
		lexer.advance()
	}

	text := lexer.sourceCode[lexer.start:lexer.current]
	if tokenType, ok := keywords[text]; ok {
		lexer.addToken(tokenType)
	} else {
		lexer.addToken(Identifier)
	}
}
