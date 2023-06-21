package lexer

import "fmt"

type TokenType int8

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal Literal
	Line    int
}

func NewToken(kind TokenType, lexeme string, literal Literal, line int) *Token {
	return &Token{
		Type:    kind,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

// TODO could be prettier
func (token *Token) String() string {
	if token.Literal != nil {
		return fmt.Sprintf("%v (value:%v)", tokenRepresentation[token.Type], token.Literal.String())
	} else {
		return fmt.Sprintf("%v", tokenRepresentation[token.Type])
	}
}

const (
	LEFT_PARANTHESIS TokenType = iota
	RIGHT_PARANTHESIS
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	DASH
	PLUS
	SEMICOLON
	SLASH
	STAR
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

var tokenRepresentation = []string{
	"LEFT_PARANTHESIS",
	"RIGHT_PARANTHESIS",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"COMMA",
	"DOT",
	"DASH",
	"PLUS",
	"SEMICOLON",
	"SLASH",
	"STAR",
	"BANG",
	"BANG_EQUAL",
	"EQUAL",
	"EQUAL_EQUAL",
	"GREATER",
	"GREATER_EQUAL",
	"LESS",
	"LESS_EQUAL",
	"IDENTIFIER",
	"STRING_LITERAL",
	"NUMBER_LITERAL",
	"AND",
	"CLASS",
	"ELSE",
	"FALSE",
	"FUN",
	"FOR",
	"IF",
	"NIL",
	"OR",
	"PRINT",
	"RETURN",
	"SUPER",
	"THIS",
	"TRUE",
	"VAR",
	"WHILE",
	"EOF",
}
