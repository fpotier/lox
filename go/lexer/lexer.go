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

func (l *Lexer) Tokens() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, *NewToken(EOF, "", nil, l.line))

	return l.tokens
}

func (l *Lexer) addToken(kind TokenType) {
	l.addTokenWithLiteral(kind, nil)
}

func (l *Lexer) addTokenWithLiteral(kind TokenType, literal Literal) {
	text := l.sourceCode[l.start:l.current]
	l.tokens = append(l.tokens, *NewToken(kind, text, literal, l.line))
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		l.line++
	case '(':
		l.addToken(LeftParenthesis)
	case ')':
		l.addToken(RightParenthesis)
	case '{':
		l.addToken(LeftBrace)
	case '}':
		l.addToken(RightBrace)
	case ';':
		l.addToken(Semicolon)
	case ',':
		l.addToken(Comma)
	case '.':
		l.addToken(Dot)
	case '-':
		l.addToken(Dash)
	case '+':
		l.addToken(Plus)
	case '*':
		l.addToken(Star)
	case '!':
		if l.match('=') {
			l.addToken(BangEqual)
		} else {
			l.addToken(Bang)
		}
	case '=':
		if l.match('=') {
			l.addToken(EqualEqual)
		} else {
			l.addToken(Equal)
		}
	case '<':
		if l.match('=') {
			l.addToken(LessEqual)
		} else {
			l.addToken(Less)
		}
	case '>':
		if l.match('=') {
			l.addToken(GreaterEqual)
		} else {
			l.addToken(Greater)
		}
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(Slash)
		}
	case '"':
		l.string()
	default:
		switch {
		case isDigit(c):
			l.number()
		case isAlpha(c):
			l.identifier()
		default:
			loxerror.Error(l.line, fmt.Sprintf("Unexpected character '%c'", c))
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

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.sourceCode)
}

func (l *Lexer) advance() byte {
	char := l.sourceCode[l.current]
	l.current++

	return char
}

func (l *Lexer) match(expected byte) bool {
	doesMatch := !l.isAtEnd() && l.sourceCode[l.current] == expected
	if doesMatch {
		l.current++
	}

	return doesMatch
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.sourceCode[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.sourceCode) {
		return 0
	}

	return l.sourceCode[l.current+1]
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}

		l.advance()
	}

	if l.isAtEnd() {
		loxerror.Error(l.line, "Unterminated string")
		return
	}

	l.advance() // the closing "
	stringValue := l.sourceCode[l.start+1 : l.current-1]
	l.addTokenWithLiteral(String, &StringLiteral{Value: stringValue})
}

// Lox number should only be of this form: (\d*)(\.?)(\d*)
func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		// Consume the '.'
		l.advance()

		for isDigit(l.peek()) {
			l.advance()
		}
	}

	floatValue, err := strconv.ParseFloat(l.sourceCode[l.start:l.current], 64)
	if err != nil {
		loxerror.Error(l.line,
			fmt.Sprintf("Error converting %v to float: %v", l.sourceCode[l.start:l.current], err))
		return
	}
	l.addTokenWithLiteral(Number, &NumberLiteral{Value: floatValue})
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := l.sourceCode[l.start:l.current]
	if tokenType, ok := keywords[text]; ok {
		l.addToken(tokenType)
	} else {
		l.addToken(Identifier)
	}
}
