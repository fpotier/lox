package loxerror

type ErrorReporter interface {
	Error(line int, message string)
}
type ParserError struct {
	Message string
}

func (e ParserError) Error() string {
	return e.Message
}

type RuntimeError struct {
	Message string
}

func (e RuntimeError) Error() string {
	return e.Message
}
