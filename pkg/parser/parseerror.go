package parser

type ParseError struct {
	line    int
	message string
}

func NewParseError(line int, message string) *ParseError {
	return &ParseError{
		line:    line,
		message: message,
	}
}

func (e *ParseError) Line() int {
	return e.line
}

func (e *ParseError) Kind() string {
	return "ParseError"
}

func (e *ParseError) Message() string {
	return e.message
}
