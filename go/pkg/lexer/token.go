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

func (t *Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%v (value:%v)", tokenRepresentation[t.Type], t.Literal.String())
	}

	return fmt.Sprintf("%v", tokenRepresentation[t.Type])
}

const (
	LeftParenthesis TokenType = iota
	RightParenthesis
	LeftBrace
	RightBrace
	Comma
	Dot
	Dash
	Plus
	Semicolon
	Slash
	Star
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

var tokenRepresentation = []string{
	"LEFT_PARENTHESIS",
	"RIGHT_PARENTHESIS",
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
